package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageRequest struct {
	Title  *string             `bson:"title" json:"title" binding:"omitempty" example:"Wedding Bliss Package"`
	Type   *models.PackageType `bson:"type" json:"type" binding:"omitempty,package_type" example:"WEDDING_BLISS"`
	Photos *[]string           `bson:"photos" json:"photos" binding:"omitempty" example:"thisisbase64image1,thisisbase64image2"`
}

type PackageResponse struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	OwnerID     primitive.ObjectID  `bson:"owner_id,omitempty" json:"ownerId" ts_type:"string" example:"12345678abcd"`
	Title       string              `form:"title" bson:"title" json:"title" binding:"required" example:"Wedding Bliss Package"`
	Type        models.PackageType  `form:"type" bson:"type" json:"type" binding:"required,package_type" example:"WEDDING_BLISS"`
	PhotoUrls   []string            `bson:"photo_urls" json:"photoUrls" example:"/package/12345678abcd_1,/package/12345678abcd_2"`
	SubPackages []models.Subpackage `bson:"sub_packages,omitempty" json:"subPackages"`
}

func (item *PackageRequest) ToModel(ownerId primitive.ObjectID) *models.Package {
	return &models.Package{
		OwnerID: ownerId,
		Title:   *item.Title,
		Type:    *item.Type,
	}
}
