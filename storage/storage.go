package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZaharBorisenko/Banking_App_Goland/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	CreateAccount(account *models.Account) error
	DeleteAccount(uuid.UUID) error
	UpdateAccount(account *models.Account) error
	GetAccounts() ([]*models.Account, error)
	GetAccountById(uuid.UUID) (*models.Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrDataNotFound        = errors.New("data rows not found")
	ErrDatabase            = errors.New("database error")
	ErrCreateDatabaseTable = errors.New("error creating table")
	ErrConnectDatabase     = errors.New("error connecting database")
)

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=Bank password=admin sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectDatabase, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectDatabase, err)
	}
	log.Printf("Connected to PostgresSQL database")
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.enableUUIDExtension(); err != nil {
		return fmt.Errorf("UUID extension disabled %w", err)
	}
	return s.CreateAccountTable()
}
func (s *PostgresStore) enableUUIDExtension() error {
	_, err := s.db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")
	if err != nil {
		return fmt.Errorf("failed to enable pgcrypto extension: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id UUID DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50),
    number serial,
    balance int,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateDatabaseTable, err)
	}
	log.Printf("Account table created or already exists")

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*models.Account, error) {
	var accounts []*models.Account
	rows, err := s.db.Query(`SELECT id, first_name, last_name, number, balance, created_at, updated_at FROM account`)
	if err != nil {
		return nil, fmt.Errorf("get accounts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		account := &models.Account{}
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%w : %v", ErrDataNotFound, err)
			}
			log.Printf("databasee error in GetAccounts %s", err)
			return nil, fmt.Errorf("%w: %v", ErrDatabase, err)

		}
		accounts = append(accounts, account)
	}
	log.Printf("Found %d accounts", len(accounts))
	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id uuid.UUID) (*models.Account, error) {
	var account = &models.Account{}
	row := s.db.QueryRow(`SELECT id, first_name, last_name, number, balance, created_at, updated_at FROM account where id = $1`, id)

	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w : %v", ErrAccountNotFound, id)
		}
		log.Printf("databasee error in GetAccountById(%s): %v\n", id, err)
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return account, nil
}

func (s *PostgresStore) CreateAccount(account *models.Account) error {
	query := `INSERT INTO 
   	account
    	(first_name,last_name, number, balance) 
	VALUES 
		($1,$2,$3,$4)`

	_, err := s.db.Exec(query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
	)

	if err != nil {
		return fmt.Errorf("failed to create account %w", err)
	}

	log.Println("Account created %w")
	return nil
}

func (s *PostgresStore) UpdateAccount(account *models.Account) error {
	query := `UPDATE account 
			  SET 
				first_name = $1, 
          		last_name = $2, 
          		number = $3, 
          		balance = $4, 
       		  WHERE id = $5 
			  RETURNING id`

	result, err := s.db.Exec(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: %v", ErrAccountNotFound, account.ID)
	}

	return nil

}

func (s *PostgresStore) DeleteAccount(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE from account where id = $1", id)
	if err != nil {
		return fmt.Errorf("deletion is not completed %w", err)
	}

	return nil
}
