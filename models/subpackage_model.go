package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subpackage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	PackageID   primitive.ObjectID `bson:"package_id,omitempty" json:"packageId" ts_type:"string" example:"12345678abcd"`
	Title       string             `bson:"title,omitempty" json:"title" example:"Wedding Bliss Package"`
	Description string             `bson:"description,omitempty" json:"description" example:"This is a package for wedding"`
	Duration    int                `bson:"duration,omitempty" json:"duration" example:"60" description:"Duration in minutes"`
	Price       int                `bson:"price,omitempty" json:"price" example:"10000"`

	IsInf              bool      `bson:"is_inf,omitempty" json:"isInf" example:"false"`
	RepeatedDay        []DayName `bson:"repeated_day,omitempty" json:"repeatedDay" binding:"day_names" ts_type:"DayName[]" example:"MON,TUE,WED"`
	availableStartTime string    `bson:"available_start_time,omitempty" json:"availableStartTime" binding:"time_format" example:"15:04"`
	availableEndTime   string    `bson:"available_end_time,omitempty" json:"availableEndTime" binding:"time_format" example:"16:27"`

	availableStartDay string `bson:"available_start_day,omitempty" json:"availableStartDay" binding:"date_format" example:"2021-01-01"`
	availableEndDay   string `bson:"available_end_day,omitempty" json:"availableEndDay" binding:"date_format" example:"2021-12-31"`
}

type DayName string

const (
	Sunday    DayName = "SUN"
	Monday    DayName = "MON"
	Tuesday   DayName = "TUE"
	Wednesday DayName = "WED"
	Thursday  DayName = "THU"
	Friday    DayName = "FRI"
	Saturday  DayName = "SAT"
)

var ValidDayNames = []struct {
	Value  DayName
	TSName string
}{
	{Sunday, string(Sunday)},
	{Monday, string(Monday)},
	{Tuesday, string(Tuesday)},
	{Wednesday, string(Wednesday)},
	{Thursday, string(Thursday)},
	{Friday, string(Friday)},
	{Saturday, string(Saturday)},
}
