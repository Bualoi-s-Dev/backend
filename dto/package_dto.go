package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
)

type PackageRequest struct {
	Title  *string             `bson:"title" json:"title" binding:"omitempty" example:"Wedding Bliss Package"`
	Type   *models.PackageType `bson:"type" json:"type" binding:"omitempty,package_type" example:"WEDDING_BLISS"`
	Photos *[]string           `bson:"photos" json:"photos" binding:"omitempty" example:"thisisbase64image1,thisisbase64image2"`
}

type PackageStrictRequest struct {
	Title  string             `bson:"title" json:"title" binding:"required" example:"Wedding Bliss Package"`
	Type   models.PackageType `bson:"type" json:"type" binding:"required,package_type" example:"WEDDING_BLISS"`
	Photos []string           `bson:"photos" json:"photos" binding:"required" example:"thisisbase64image1,thisisbase64image2"`
}
