package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/lib/pq"
)

type healthRecordRepository struct {
	db *sql.DB
}

func NewHealthRecordRepository(db *sql.DB) *healthRecordRepository {
	return &healthRecordRepository{db: db}
}

func (hr *healthRecordRepository) Create(ctx context.Context, healthRecord *domain.HealthRecord) (string, error) {
	var id string
	query := `
		INSERT INTO health_records 
			(user_id , health_record_type_id, value , notes ,   updated_at, created_at ,recorded_at)
		VALUES 
			($1,
			$2,
			$3,
			$4, 
			now(),
			now(),
			$5)
		RETURNING ID
	`
	err := hr.db.QueryRowContext(ctx,
		query,
		healthRecord.UserID,
		healthRecord.HealthRecordTypeID,
		healthRecord.Value,
		healthRecord.Notes,
		healthRecord.RecordedAt,
	).Scan(
		&id,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		pqErr, ok := err.(*pq.Error)

		if ok && pqErr.Code == "23503" && pqErr.Constraint == "health_records" {
			return id, domain.ErrInvalidUserOrHealthRecordType
		}

		return id, err
	}

	return id, nil
}

func (hr *healthRecordRepository) GetByID(ctx context.Context, id string) (*domain.HealthRecord, error) {
	healthRecord := &domain.HealthRecord{}
	query := `
	SELECT 
		id,
		user_id, 
		health_record_type_id, 
		value,
		notes, 
		created_at,
		updated_at,
		recorded_at  
	FROM health_records
	WHERE id = $1
	`
	err := hr.db.QueryRowContext(ctx,
		query,
		id).Scan(&healthRecord.ID,
		&healthRecord.UserID,
		&healthRecord.HealthRecordTypeID,
		&healthRecord.Value,
		&healthRecord.Notes,
		&healthRecord.CreatedAt,
		&healthRecord.UpdatedAt,
		&healthRecord.RecordedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrHealthRecordNotFound
		}
		return nil, err
	}

	return healthRecord, err
}

func (hr *healthRecordRepository) ListByUserID(ctx context.Context, userId string) ([]*domain.HealthRecord, error) {
	query := `
	SELECT 
		id,
		user_id, 
		health_record_type_id, 
		value,
		notes  
	FROM health_records
	WHERE user_id = $1
	`
	rows, err := hr.db.QueryContext(ctx, query, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrHealthRecordNotFound
		}
		return nil, err
	}

	var healthRecords []*domain.HealthRecord
	for rows.Next() {
		var healthRecord domain.HealthRecord
		rows.Scan(
			&healthRecord.ID,
			&healthRecord.UserID,
			&healthRecord.HealthRecordTypeID,
			&healthRecord.Value,
			&healthRecord.Notes,
		)
		healthRecords = append(healthRecords, &healthRecord)
	}

	return healthRecords, nil
}

func (hr *healthRecordRepository) Update(ctx context.Context, id string, healthRecord *domain.HealthRecord) error {
	query := `
		UPDATE health_records
		SET 
			VALUE = $2,
			NOTES = $3,
			recorded_at = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	result, err := hr.db.ExecContext(ctx, query, id, healthRecord.Value, healthRecord.Notes, healthRecord.RecordedAt)

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrHealthRecordNotFound
	}

	return nil
}

func (hr *healthRecordRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM health_records 
		WHERE id = $1;
	`
	result, err := hr.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrHealthRecordNotFound
	}

	return nil
}
