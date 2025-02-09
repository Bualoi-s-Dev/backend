package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string"`
	Email    string             `bson:"email" json:"email"`
	Name     string             `bson:"name,omitempty" json:"name"`
	Gender   string             `bson:"gender,omitempty" json:"gender"`
	Profile  string             `bson:"profile,omitempty" json:"profile"`
	Phone    string             `bson:"phone,omitempty" json:"phone"`
	Location string             `bson:"location,omitempty" json:"location"`

	//Photographer Info
	IsPhotographer   bool                 `bson:"is_photographer" json:"isPhotographer"`
	BankName         BankName             `bson:"bank_name,omitempty" json:"bankName"`
	BankAccount      string               `bson:"bank_account,omitempty" json:"bankAccount"`
	LineID           string               `bson:"line_id,omitempty" json:"lineID"`
	Facebook         string               `bson:"facebook,omitempty" json:"facebook"`
	Instagram        string               `bson:"instagram,omitempty" json:"instagram"`
	ShowcasePackages []primitive.ObjectID `bson:"showcase_packages,omitempty" json:"showcasePackages" ts_type:"string"`
	Packages         []primitive.ObjectID `bson:"photographer_packages,omitempty" json:"packages" ts_type:"string"`
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
