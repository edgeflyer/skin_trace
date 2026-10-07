package record

import (
	"path/filepath"
	"time"
)

type Response struct {
	ID        uint64    `json:"id"`
	PhotoURL  string    `json:"photo_url"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

func toResponse(rec Record) Response {
	return Response{
		ID:        rec.ID,
		PhotoURL:  "/uploads/" + filepath.Base(rec.PhotoPath),
		Note:      rec.Note,
		CreatedAt: rec.CreatedAt,
	}
}
