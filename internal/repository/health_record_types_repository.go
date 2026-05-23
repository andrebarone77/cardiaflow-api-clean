package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/lib/pq"
)

type healthRecordTypeRepository struct {
	db *sql.DB
}

func NewHealthRecordTypeRepository(db *sql.DB) *healthRecordTypeRepository {
	return &healthRecordTypeRepository{db: db}
}

func (r *healthRecordTypeRepository) Create(ctx context.Context, h *domain.HealthRecordType) (string, error) {
	query := `
		INSERT INTO health_record_types 
		(name, code, unit)
		VALUES
		($1, $2, $3)
		RETURNING id
	`
	var id string

	err := r.db.QueryRowContext(ctx,
		query,
		h.Name,
		h.Code,
		h.Unit,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating health record type: %v", err)
		pqErr, ok := err.(*pq.Error)

		if ok && pqErr.Code == "23505" && pqErr.Constraint == "health_record_types_code_key" {
			return "", domain.ErrHealthRecordTypeAlreadyExists
		}

		return "", err
	}

	return id, nil

}

func (r *healthRecordTypeRepository) GetById(ctx context.Context, id string) (*domain.HealthRecordType, error) {
	healthRecordType := &domain.HealthRecordType{}
	query := `
			SELECT id, name, code 
			FROM health_record_types
			WHERE id = $1
			`

	err := r.db.QueryRowContext(ctx,
		query,
		id).Scan(&healthRecordType.ID,
		&healthRecordType.Name,
		&healthRecordType.Code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrHealthRecordTypeNotFound
		}
		return nil, err
	}

	return healthRecordType, nil
}

func (r *healthRecordTypeRepository) IsSystem(ctx context.Context, id string) (bool, error) {
	var result bool
	query := `
		SELECT is_system 
		FROM health_record_types
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx,
		query,
		id).Scan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrHealthRecordTypeNotFound
		}
		return false, err
	}

	return result, nil
}

func (r *healthRecordTypeRepository) GetByCode(ctx context.Context, code string) (*domain.HealthRecordType, error) {
	healthRecordType := &domain.HealthRecordType{}
	query := `
			SELECT id, name, code 
			FROM health_record_types
			WHERE code = $1
			`

	err := r.db.QueryRowContext(ctx,
		query,
		code).Scan(&healthRecordType.ID,
		&healthRecordType.Name,
		&healthRecordType.Code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrHealthRecordTypeNotFound
		}
		return nil, err
	}

	return healthRecordType, nil
}

func (r *healthRecordTypeRepository) GetAll(ctx context.Context) ([]*domain.HealthRecordType, error) {
	query := `
		SELECT id, name, code, is_system 
		FROM health_record_types
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrHealthRecordTypeNotFound
		}
		return nil, err
	}

	var healthRecordTypes []*domain.HealthRecordType

	for rows.Next() {
		var healthRecordType domain.HealthRecordType

		err := rows.Scan(
			&healthRecordType.ID,
			&healthRecordType.Name,
			&healthRecordType.Code,
			&healthRecordType.IsSystem,
		)
		if err != nil {
			return nil, err
		}

		healthRecordTypes = append(healthRecordTypes, &healthRecordType)
	}

	return healthRecordTypes, nil
}

func (r *healthRecordTypeRepository) Update(ctx context.Context, id string, h *domain.HealthRecordType) error {
	query := `
		UPDATE health_record_types
		SET name = $1,
		    code = $2,
		    unit = $3,
		    updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx,
		query,
		h.Name,
		h.Code,
		h.Unit,
		id,
	)

	if err != nil {
		log.Printf("Error updating health_record_type: %v", err)

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "health_record_types_code_key" {
				return domain.ErrCodeAlreadyExists
			}
		}

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrHealthRecordTypeNotFound
	}

	return nil
}
func (r *healthRecordTypeRepository) Delete(ctx context.Context, id string) error {

	query := `
		DELETE FROM health_record_types
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrHealthRecordTypeNotFound
	}

	return nil
}
