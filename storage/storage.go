package storage

import (
	"database/sql"
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

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=Bank password=admin sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected to PostgreSQL database")
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.enableUUIDExtension(); err != nil {
		return err
	}
	return s.CreateAccountTable()
}
func (s *PostgresStore) enableUUIDExtension() error {
	_, err := s.db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")
	return err
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id UUID DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50),
    number serial,
    balance int,
    created_at timestamp
)`

	_, err := s.db.Exec(query)
	if err != nil {
		fmt.Errorf("failed to create account table %w", err)
	}
	log.Println("Account table created or already exists")

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*models.Account, error) {
	var accounts []*models.Account
	rows, err := s.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, fmt.Errorf("failed to select accounts: %w", err)
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
		); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Found %d accounts", len(accounts))
	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id uuid.UUID) (*models.Account, error) {
	var account = &models.Account{}
	row := s.db.QueryRow(`SELECT * FROM account where id = $1`, id)

	if err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failde to scan account %w", err)
	}

	return account, nil
}

func (s *PostgresStore) CreateAccount(account *models.Account) error {
	query := `INSERT INTO 
   	account
    	(first_name,last_name, number, balance, created_at) 
	VALUES 
		($1,$2,$3,$4,$5)`

	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)

	if err != nil {
		fmt.Errorf("failed to create account %w", err)
	}

	log.Println("Account created %w")
	return nil
}

func (s *PostgresStore) UpdateAccount(account *models.Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE from account where id = $1", id)
	if err != nil {
		return fmt.Errorf("deletion is not completed %w", err)
	}

	return nil
}
