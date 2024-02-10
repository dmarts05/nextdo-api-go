package customer

import (
	"context"
	"errors"

	"github.com/dmarts05/nextdo-api-go/internal/shared/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, customer Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (Customer, error)
	GetByEmail(ctx context.Context, email string) (Customer, error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepository struct {
	Db *pgxpool.Pool
}

func (r PostgresRepository) Create(ctx context.Context, customer Customer) error {
	// Check if a customer with the same email already exists
	_, err := r.GetByEmail(ctx, customer.Email)
	if err == nil {
		return repository.ErrAlreadyExists{}
	}

	commandTag, err := r.Db.Exec(ctx, "INSERT INTO customer (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)", customer.FirstName, customer.LastName, customer.Email, customer.Password)

	switch {
	case err != nil:
		return repository.ErrFailedToCreate{Err: err}
	case commandTag.RowsAffected() != 1:
		return repository.ErrFailedToCreate{Err: errors.New("no rows affected")}
	}

	return nil
}

func (r PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (Customer, error) {
	var customer Customer
	err := r.Db.QueryRow(ctx, "SELECT id, first_name, last_name, email, password, created_at, updated_at FROM customer WHERE id = $1", id).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Password, &customer.CreatedAt, &customer.UpdatedAt)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return Customer{}, repository.ErrNotFound{}
	case err != nil:
		return Customer{}, repository.ErrFailedToGet{Err: err}
	}

	return customer, nil
}

func (r PostgresRepository) GetByEmail(ctx context.Context, email string) (Customer, error) {
	var customer Customer
	err := r.Db.QueryRow(ctx, "SELECT id, first_name, last_name, email, password, created_at, updated_at FROM customer WHERE email = $1", email).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Password, &customer.CreatedAt, &customer.UpdatedAt)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return Customer{}, repository.ErrNotFound{}
	case err != nil:
		return Customer{}, repository.ErrFailedToGet{Err: err}
	}

	return customer, nil
}

func (r PostgresRepository) Update(ctx context.Context, customer Customer) error {
	// Email cannot be updated
	commandTag, err := r.Db.Exec(ctx, "UPDATE customer SET first_name = $1, last_name = $2, password = $3, updated_at = now() WHERE id = $4", customer.FirstName, customer.LastName, customer.Password, customer.ID)

	switch {
	case err != nil:
		return repository.ErrFailedToUpdate{Err: err}
	case commandTag.RowsAffected() != 1:
		return repository.ErrFailedToUpdate{Err: errors.New("no rows affected")}
	}

	return nil
}

func (r PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	commandTag, err := r.Db.Exec(ctx, "DELETE FROM customer WHERE id = $1", id)

	switch {
	case err != nil:
		return repository.ErrFailedToDelete{Err: err}
	case commandTag.RowsAffected() != 1:
		return repository.ErrFailedToDelete{Err: errors.New("no rows affected")}
	}

	return nil
}
