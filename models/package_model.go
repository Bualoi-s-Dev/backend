package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Package struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title  string             `bson:"title" json:"title" binding:"required"`
	Type   PackageType        `bson:"type" json:"type" binding:"required,package_type"`
	Photos []string           `bson:"photos" json:"photos"`
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

var validPackageTypes = []PackageType{
	WeddingBliss,
	BirthdayShoots,
	MaternityGlow,
	FamilyPortraits,
	GraduationMemories,
	NewbornMoments,
	EngagementLoveStory,
	CorporateHeadshots,
	FashionEditorial,
	TravelDiaries,
	Other,
}
