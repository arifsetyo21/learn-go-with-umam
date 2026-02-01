package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.ProductCategoryView, error) {
	query := "SELECT prod.id, prod.name, prod.price, prod.stock, cat.name FROM products prod JOIN categories cat ON prod.category_id = cat.id"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productCategory := make([]models.ProductCategoryView, 0)
	for rows.Next() {
		var p models.ProductCategoryView
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryName)

		if err != nil {
			return nil, err
		}
		productCategory = append(productCategory, p)
	}

	return productCategory, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetByID(id int) (*models.ProductCategoryView, error) {
	query := "SELECT prod.id, prod.name, prod.price, prod.stock, cat.name FROM products prod JOIN categories cat ON prod.category_id = cat.id WHERE prod.id = $1"

	var p models.ProductCategoryView
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Product tidak ditemukan")
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product tidak ditemukan")
	}

	return err
}
