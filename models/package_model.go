package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Package struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	OwnerID   primitive.ObjectID `bson:"owner_id,omitempty" json:"ownerId" ts_type:"string" example:"12345678abcd"`
	Title     string             `form:"title" bson:"title" json:"title" binding:"required" example:"Wedding Bliss Package"`
	Type      PackageType        `form:"type" bson:"type" json:"type" binding:"required,package_type" example:"WEDDING_BLISS"`
	PhotoUrls []string           `bson:"photo_urls" json:"photoUrls" example:"/package/12345678abcd_1,/package/12345678abcd_2"`
}

type PackageType string

const (
	WeddingBliss        PackageType = "WEDDING_BLISS"
	BirthdayShoots      PackageType = "BIRTHDAY_SHOOTS"
	MaternityGlow       PackageType = "MATERNITY_GLOW"
	FamilyPortraits     PackageType = "FAMILY_PORTRAITS"
	GraduationMemories  PackageType = "GRADUATION_MEMORIES"
	NewbornMoments      PackageType = "NEWBORN_MOMENTS"
	EngagementLoveStory PackageType = "ENGAGEMENT_LOVE_STORY"
	CorporateHeadshots  PackageType = "CORPORATE_HEADSHOTS"
	FashionEditorial    PackageType = "FASHION_EDITORIAL"
	TravelDiaries       PackageType = "TRAVEL_DIARIES"
	Other               PackageType = "OTHER"
)

var ValidPackageTypes = []struct {
	Value  PackageType
	TSName string
}{
	{WeddingBliss, string(WeddingBliss)},
	{BirthdayShoots, string(BirthdayShoots)},
	{MaternityGlow, string(MaternityGlow)},
	{FamilyPortraits, string(FamilyPortraits)},
	{GraduationMemories, string(GraduationMemories)},
	{NewbornMoments, string(NewbornMoments)},
	{EngagementLoveStory, string(EngagementLoveStory)},
	{CorporateHeadshots, string(CorporateHeadshots)},
	{FashionEditorial, string(FashionEditorial)},
	{TravelDiaries, string(TravelDiaries)},
	{Other, string(Other)},
}

type Subpackage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd" binding:"required"`
	PackageID   primitive.ObjectID `bson:"package_id,omitempty" json:"packageId" ts_type:"string" example:"12345678abcd" binding:"required"`
	Name        string             `form:"name" bson:"name" json:"name" binding:"required" example:"Basic Package"`
	Description string             `form:"description" bson:"description" json:"description" binding:"required" example:"Basic package for wedding bliss"`
	Price       int                `form:"price" bson:"price" json:"price" binding:"required" example:"5000"`         // (THB)
	Duration    int                `form:"duration" bson:"duration" json:"duration" binding:"required" example:"120"` // (minutes)
}
