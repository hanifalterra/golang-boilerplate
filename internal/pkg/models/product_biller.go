package models

import (
	"time"
)

type ProductBiller struct {
	ID        uint
	ProductID uint
	BillerID  uint
	IsActive  bool
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
	DeletedBy string
}

func (pb *ProductBiller) ToResponse() *ProductBillerResponse {
	return &ProductBillerResponse{
		ID:        pb.ID,
		ProductID: pb.ProductID,
		BillerID:  pb.BillerID,
		IsActive:  pb.IsActive,
		CreatedAt: pb.CreatedAt,
		CreatedBy: pb.CreatedBy,
		UpdatedAt: pb.UpdatedAt,
		UpdatedBy: pb.UpdatedBy,
		DeletedAt: pb.DeletedAt,
		DeletedBy: pb.DeletedBy,
	}
}

type CreateProductBillerRequest struct {
	ProductID uint `json:"product_id"`
	BillerID  uint `json:"biller_id"`
	IsActive  bool `json:"is_active"`
}

func (pb *CreateProductBillerRequest) ToDomain() *ProductBiller {
	return &ProductBiller{
		ProductID: pb.ProductID,
		BillerID:  pb.BillerID,
		IsActive:  pb.IsActive,
	}
}

type UpdateProductBillerRequest struct {
	IsActive bool `json:"is_active"`
}

func (pb *UpdateProductBillerRequest) ToDomain() *ProductBiller {
	return &ProductBiller{
		IsActive: pb.IsActive,
	}
}

type ProductBillerResponse struct {
	ID        uint       `json:"id"`
	ProductID uint       `json:"product_id"`
	BillerID  uint       `json:"biller_id"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy string     `json:"deleted_by"`
}
