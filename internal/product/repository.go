package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// Repository defines the interface for product data access
type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	List(ctx context.Context, page, pageSize int) ([]Product, int, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
}

// MySQLRepository implements Repository using MySQL
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new MySQL product repository
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// Create inserts a new product
func (r *MySQLRepository) Create(ctx context.Context, p *Product) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	query := `
		INSERT INTO products (id, sku, name, description, price_cents, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.SKU, p.Name, p.Description, p.PriceCents, p.Category, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		// Check for duplicate key error (MySQL error 1062)
		if isDuplicateKeyError(err) {
			return ErrSKUAlreadyExists
		}
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *MySQLRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT id, sku, name, description, price_cents, category, created_at, updated_at
		FROM products
		WHERE id = ?
	`
	p := &Product{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &desc, &p.PriceCents, &p.Category, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if desc.Valid {
		p.Description = &desc.String
	}
	return p, nil
}

// GetBySKU retrieves a product by SKU
func (r *MySQLRepository) GetBySKU(ctx context.Context, sku string) (*Product, error) {
	query := `
		SELECT id, sku, name, description, price_cents, category, created_at, updated_at
		FROM products
		WHERE sku = ?
	`
	p := &Product{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&p.ID, &p.SKU, &p.Name, &desc, &p.PriceCents, &p.Category, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	if desc.Valid {
		p.Description = &desc.String
	}
	return p, nil
}

// List retrieves products with pagination
func (r *MySQLRepository) List(ctx context.Context, page, pageSize int) ([]Product, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM products`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, sku, name, description, price_cents, category, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var desc sql.NullString
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &desc, &p.PriceCents, &p.Category, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		if desc.Valid {
			p.Description = &desc.String
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return products, total, nil
}

// Update updates a product
func (r *MySQLRepository) Update(ctx context.Context, p *Product) error {
	p.UpdatedAt = time.Now()
	query := `
		UPDATE products
		SET sku = ?, name = ?, description = ?, price_cents = ?, category = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		p.SKU, p.Name, p.Description, p.PriceCents, p.Category, p.UpdatedAt, p.ID,
	)
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrSKUAlreadyExists
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

// Delete removes a product
func (r *MySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

// isDuplicateKeyError checks if the error is a MySQL duplicate key error (1062)
func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}