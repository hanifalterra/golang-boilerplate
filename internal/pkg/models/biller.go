package models

import (
	"time"
)

type Biller struct {
	ID        uint
	Label     string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
	DeletedBy string
}

func (b *Biller) ToResponse() *BillerResponse {
	return &BillerResponse{
		ID:        b.ID,
		Label:     b.Label,
		CreatedAt: b.CreatedAt,
		CreatedBy: b.CreatedBy,
		UpdatedAt: b.UpdatedAt,
		UpdatedBy: b.UpdatedBy,
		DeletedAt: b.DeletedAt,
		DeletedBy: b.DeletedBy,
	}
}

type CreateBillerRequest struct {
	Label string `json:"label"`
}

func (b *CreateBillerRequest) ToDomain() *Biller {
	return &Biller{
		Label: b.Label,
	}
}

type UpdateBillerRequest struct {
	Label string `json:"label"`
}

func (b *UpdateBillerRequest) ToDomain() *Biller {
	return &Biller{
		Label: b.Label,
	}
}

type BillerResponse struct {
	ID        uint       `json:"id"`
	Label     string     `json:"label"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy string     `json:"deleted_by"`
}
