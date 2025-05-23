package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	Email    string             `bson:"email" json:"email" example:"ceo.meen@gmail.com"`
	Name     string             `bson:"name,omitempty" json:"name" example:"Meen"`
	Gender   string             `bson:"gender,omitempty" json:"gender" example:"LGTV"`
	Profile  string             `bson:"profile,omitempty" json:"profile" example:"/profile/12345678abcd"`
	Phone    string             `bson:"phone,omitempty" json:"phone" example:"0812345678"`
	Location string             `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`

	//Photographer Info
	Role             UserRole             `bson:"role,omitempty" json:"role" binding:"omitempty,user_role" example:"Photographer"`
	Description      string               `bson:"description,omitempty" json:"description" example:"I'm a photographer"`
	BankName         BankName             `bson:"bank_name,omitempty" json:"bankName" binding:"omitempty,bank_name" example:"KRUNG_THAI_BANK"`
	BankAccount      string               `bson:"bank_account,omitempty" json:"bankAccount" example:"1234567890"`
	LineID           string               `bson:"line_id,omitempty" json:"lineID" example:"@meen"`
	Facebook         string               `bson:"facebook,omitempty" json:"facebook" example:"Meen"`
	Instagram        string               `bson:"instagram,omitempty" json:"instagram" example:"Meen"`
	ShowcasePackages []primitive.ObjectID `bson:"showcase_packages,omitempty" json:"showcasePackages" ts_type:"string[]" example:"12345678abcd,12345678abcd"`

	// Payment Info
	StripeCustomerID *string `bson:"stripe_customer_id,omitempty" json:"stripeCustomerId" example:"12345678abcd"`
	StripeAccountID  *string `bson:"stripe_account_id,omitempty" json:"stripeAccountId" example:"12345678abcd"`
}

func NewUser(email string) *User {
	return &User{
		ID:               primitive.NewObjectID(),
		Email:            email,
		Name:             "",
		Gender:           "",
		Profile:          "",
		Phone:            "",
		Location:         "",
		Role:             Guest,
		ShowcasePackages: []primitive.ObjectID{},
	}
}

type UserRole string

const (
	Photographer UserRole = "Photographer"
	Customer     UserRole = "Customer"
	Guest        UserRole = "Guest"
)

var ValidUserRoles = []struct {
	Value  UserRole
	TSName string
}{
	{Photographer, string(Photographer)},
	{Customer, string(Customer)},
	{Guest, string(Guest)},
}

type BankName string

const (
	KrungThaiBank         BankName = "KRUNG_THAI_BANK"
	BangkokBank           BankName = "BANGKOK_BANK"
	SiamCommercialBank    BankName = "SIAM_COMMERCIAL_BANK"
	KasikornBank          BankName = "KASIKORN_BANK"
	TMBThanachartBank     BankName = "TMB_THANACHART_BANK"
	KrungsriBank          BankName = "KRUNGSRI_BANK"
	GovernmentSavingsBank BankName = "GOVERNMENT_SAVINGS_BANK"
	ThaiMilitaryBank      BankName = "THAI_MILITARY_BANK"
	UOBThailand           BankName = "UOB_THAILAND"
	CIMBThailand          BankName = "CIMB_THAILAND"
	StandardChartered     BankName = "STANDARD_CHARTERED"
	ICBCThailand          BankName = "ICBC_THAILAND"
)

var ValidBankNames = []struct {
	Value  BankName
	TSName string
}{
	{KrungThaiBank, string(KrungThaiBank)},
	{BangkokBank, string(BangkokBank)},
	{SiamCommercialBank, string(SiamCommercialBank)},
	{KasikornBank, string(KasikornBank)},
	{TMBThanachartBank, string(TMBThanachartBank)},
	{KrungsriBank, string(KrungsriBank)},
	{GovernmentSavingsBank, string(GovernmentSavingsBank)},
	{ThaiMilitaryBank, string(ThaiMilitaryBank)},
	{UOBThailand, string(UOBThailand)},
	{CIMBThailand, string(CIMBThailand)},
	{StandardChartered, string(StandardChartered)},
	{ICBCThailand, string(ICBCThailand)},
}
