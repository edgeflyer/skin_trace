package record

import "time"

type Record struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	PhotoPath string    `json:"photo_path" gorm:"size:512;not null"`
	Note      string    `json:"note" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

func (Record) TableName() string {
	return "skin_records"
}
