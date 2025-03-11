package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRequest struct {
	Name     *string `bson:"name,omitempty" json:"name" example:"Meen"`
	Gender   *string `bson:"gender,omitempty" json:"gender" example:"LGTV"`
	Profile  *string `bson:"profile,omitempty" json:"profile" example:"base64123123123"`
	Phone    *string `bson:"phone,omitempty" json:"phone" example:"0812345678"`
	Location *string `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`

	Role             *models.UserRole      `bson:"role,omitempty" json:"role" binding:"omitempty,user_role" example:"Photographer"`
	Description      *string               `bson:"description,omitempty" json:"description" example:"I'm a photographer"`
	BankName         *models.BankName      `bson:"bank_name,omitempty" json:"bankName" example:"KRUNG_THAI_BANK"`
	BankAccount      *string               `bson:"bank_account,omitempty" json:"bankAccount" example:"1234567890"`
	LineID           *string               `bson:"line_id,omitempty" json:"lineID" example:"@meen"`
	Facebook         *string               `bson:"facebook,omitempty" json:"facebook" example:"Meen"`
	Instagram        *string               `bson:"instagram,omitempty" json:"instagram" example:"Meen"`
	ShowcasePackages *[]primitive.ObjectID `bson:"showcase_packages,omitempty" json:"showcasePackages" ts_type:"string[]" example:"12345678abcd,12345678abcd"`
}

type UserResponse struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	Email    string             `bson:"email" json:"email" example:"ceo.meen@gmail.com"`
	Name     string             `bson:"name,omitempty" json:"name" example:"Meen"`
	Gender   string             `bson:"gender,omitempty" json:"gender" example:"LGTV"`
	Profile  string             `bson:"profile,omitempty" json:"profile" example:"/profile/12345678abcd"`
	Phone    string             `bson:"phone,omitempty" json:"phone" example:"0812345678"`
	Location string             `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`

	Role             models.UserRole   `bson:"role,omitempty" json:"role" binding:"omitempty,user_role" example:"Photographer"`
	Description      string            `bson:"description,omitempty" json:"description" example:"I'm a photographer"`
	BankName         models.BankName   `bson:"bank_name,omitempty" json:"bankName" example:"KRUNG_THAI_BANK"`
	BankAccount      string            `bson:"bank_account,omitempty" json:"bankAccount" example:"1234567890"`
	LineID           string            `bson:"line_id,omitempty" json:"lineID" example:"@meen"`
	Facebook         string            `bson:"facebook,omitempty" json:"facebook" example:"Meen"`
	Instagram        string            `bson:"instagram,omitempty" json:"instagram" example:"Meen"`
	ShowcasePackages []PackageResponse `bson:"showcase_packages,omitempty" json:"showcasePackages" ts_type:"PackageResponse[]"`
	Packages         []PackageResponse `bson:"photographer_packages,omitempty" json:"photographerPackages" ts_type:"PackageResponse[]"`
}

type AuthUserCredentials struct {
	Email    string `json:"email" binding:"required" example:"meen@gmail.com"`
	Password string `json:"password" binding:"required" example:"12345678"`
}
