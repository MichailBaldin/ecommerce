package repository

import (
	"database/sql"
	"ecommerce/services/users/models"
	"time"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(databaseURL string) (UserRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &PostgresRepo{db: db}
	if err := repo.createTables(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *PostgresRepo) Create(user *models.User) error {
	query := `INSERT INTO users (name, email, created_at, updated_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.QueryRow(query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}

func (r *PostgresRepo) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // возвращаем nil, nil если пользователь не найден
	}

	return user, err
}

func (r *PostgresRepo) createTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    )`

	_, err := r.db.Exec(query)
	return err
}

func (r *PostgresRepo) Close() error {
	return r.db.Close()
}
