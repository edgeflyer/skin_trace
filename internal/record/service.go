package record

import (
	"context"
	"errors"
	"io"

	"skintrace/internal/storage"
)

type Service struct {
	repo  *Repository
	files *storage.Local
}

func NewService(repo *Repository, files *storage.Local) *Service {
	return &Service{
		repo:  repo,
		files: files,
	}
}

func (s *Service) Create(
	ctx context.Context,
	photo io.ReadSeeker,
	note string,
) (*Record, error) {
	ext, err := detectPhotoExtension(photo)
	if err != nil {
		return nil, err
	}

	path, err := s.files.Save(photo, ext)
	if err != nil {
		return nil, err
	}

	rec := &Record{
		PhotoPath: path,
		Note:      note,
	}

	if err := s.repo.Create(ctx, rec); err != nil {
		cleanupErr := s.files.Delete(path)
		return nil, errors.Join(err, cleanupErr)
	}

	return rec, nil
}

func (s *Service) List(ctx context.Context) ([]Record, error) {
	return s.repo.List(ctx)
}

func (s *Service) Recent(
	ctx context.Context,
	limit int,
) ([]Record, error) {
	if limit < 1 || limit > 30 {
		return nil, errors.New("limit must be between 1 and 30")
	}

	return s.repo.Recent(ctx, limit)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uint64,
) (*Record, error) {
	if id == 0 {
		return nil, errors.New("record_id must be greater than 0")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) ReadPhoto(
	ctx context.Context,
	id uint64,
) ([]byte, error) {
	rec, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.files.Read(rec.PhotoPath)
}
