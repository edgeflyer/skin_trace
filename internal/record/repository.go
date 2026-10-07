package record

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("skin record not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, rec *Record) error {
	if err := r.db.WithContext(ctx).Create(rec).Error; err != nil {
		return fmt.Errorf("create skin record: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]Record, error) {
	records := make([]Record, 0)

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Order("id DESC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("list skin records: %w", err)
	}

	return records, nil
}

func (r *Repository) Recent(
	ctx context.Context,
	limit int,
) ([]Record, error) {
	records := make([]Record, 0)

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Order("id DESC").
		Limit(limit).
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("query recent skin records: %w", err)
	}

	return records, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uint64,
) (*Record, error) {
	var rec Record

	err := r.db.WithContext(ctx).First(&rec, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get skin record: %w", err)
	}

	return &rec, nil
}
