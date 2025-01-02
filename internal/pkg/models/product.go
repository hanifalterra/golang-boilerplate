package models

import (
	"time"
)

type Product struct {
	ID        uint
	Label     string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
	DeletedBy string
}

func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:        p.ID,
		Label:     p.Label,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy,
		DeletedAt: p.DeletedAt,
		DeletedBy: p.DeletedBy,
	}
}

type CreateProductRequest struct {
	Label string `json:"label"`
}

func (p *CreateProductRequest) ToEntity() *Product {
	return &Product{
		Label: p.Label,
	}
}

type UpdateProductRequest struct {
	Label string `json:"label"`
}

func (p *UpdateProductRequest) ToEntity() *Product {
	return &Product{
		Label: p.Label,
	}
}

type ProductResponse struct {
	ID        uint       `json:"id"`
	Label     string     `json:"label"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy string     `json:"deleted_by"`
}
