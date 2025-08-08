package repository

import (
	"database/sql"
	"ecommerce/services/products/models"
	"time"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(databaseURL string) (ProductRepository, error) {
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

func (r *PostgresRepo) Create(product *models.Product) error {
	query := `INSERT INTO products (name, description, price, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	return r.db.QueryRow(query, product.Name, product.Description, product.Price, product.CreatedAt, product.UpdatedAt).Scan(&product.ID)
}

func (r *PostgresRepo) GetByID(id int) (*models.Product, error) {
	product := &models.Product{}
	query := `SELECT id, name, description, price, created_at, updated_at FROM products WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // возвращаем nil, nil если товар не найден
	}

	return product, err
}

func (r *PostgresRepo) createTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(200) NOT NULL,
        description TEXT,
        price DECIMAL(10,2) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    )`

	_, err := r.db.Exec(query)
	return err
}

func (r *PostgresRepo) Close() error {
	return r.db.Close()
}
