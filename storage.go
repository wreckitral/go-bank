package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connString := "user=postgres dbname=postgres password=goprojectbank sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		firstName varchar(50),
		lastName varchar(50),
		number serial,
		encryptedPassword varchar(100),
		balance serial,
		createdAt timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account (firstName, lastName, number, encryptedPassword, balance, createdAt)	values ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.EncryptedPassword, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	row, err := s.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		return scanIntoAccount(row)
	}
	return nil, fmt.Errorf("account with number [%d] not found", number)
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	row, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		return scanIntoAccount(row)
	}
	return nil, fmt.Errorf("account with id [%d] not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := Account{}

	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.EncryptedPassword, &account.Balance, &account.CreatedAt)

	return &account, err

}
