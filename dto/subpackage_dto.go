package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubpackageRequest struct {
	Title       *string `bson:"title" json:"title" binding:"omitempty" example:"Wedding Bliss Package"`
	Description *string `bson:"description" json:"description" binding:"omitempty" example:"This is a package for wedding"`
	Price       *int    `bson:"price" json:"price" binding:"omitempty" example:"10000"`
	Duration    *int    `bson:"duration" json:"duration" binding:"omitempty" example:"60" description:"Duration in minutes"`

	IsInf *bool `bson:"is_inf" json:"isInf" binding:"omitempty,isInf_rule" example:"false"`
	// TODO: change tsgen type to Dayname
	RepeatedDay        *[]models.DayName `bson:"repeated_day" json:"repeatedDay" binding:"omitempty,day_names" example:"MON,TUE,WED"`
	AvailableStartTime *string           `bson:"available_start_time" json:"availableStartTime" binding:"omitempty,time_format" example:"15:04"`
	AvailableEndTime   *string           `bson:"available_end_time" json:"availableEndTime" binding:"omitempty,time_format" example:"16:27"`
	AvailableStartDay  *string           `bson:"available_start_day" json:"availableStartDay" binding:"omitempty,date_format" example:"2021-01-01"`
	AvailableEndDay    *string           `bson:"available_end_day" json:"availableEndDay" binding:"omitempty,date_format" example:"2021-12-31"`
}

type SubpackageResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	PackageID   primitive.ObjectID `bson:"package_id,omitempty" json:"packageId" ts_type:"string" example:"12345678abcd"`
	Title       string             `bson:"title" json:"title" binding:"omitempty" example:"Wedding Bliss Package"`
	Description string             `bson:"description" json:"description" binding:"omitempty" example:"This is a package for wedding"`
	Price       int                `bson:"price" json:"price" binding:"omitempty" example:"10000"`
	Duration    int                `bson:"duration" json:"duration" binding:"omitempty" example:"60" description:"Duration in minutes"`

	IsInf bool `bson:"is_inf" json:"isInf" binding:"omitempty,isInf_rule" example:"false"`
	// TODO: change tsgen type to Dayname
	RepeatedDay        []models.DayName `bson:"repeated_day" json:"repeatedDay" binding:"omitempty,day_names" example:"MON,TUE,WED"`
	AvailableStartTime string           `bson:"available_start_time" json:"availableStartTime" binding:"omitempty,time_format" example:"15:04"`
	AvailableEndTime   string           `bson:"available_end_time" json:"availableEndTime" binding:"omitempty,time_format" example:"16:27"`
	AvailableStartDay  string           `bson:"available_start_day" json:"availableStartDay" binding:"omitempty,date_format" example:"2021-01-01"`
	AvailableEndDay    string           `bson:"available_end_day" json:"availableEndDay" binding:"omitempty,date_format" example:"2021-12-31"`

	BusyTimes []models.BusyTime `bson:"busy_times" json:"busyTimes" binding:"omitempty"`

	// TODO: remove this field
	BusyTimeMap map[string][]models.BusyTime `bson:"busy_time_map" json:"busyTimeMap" binding:"omitempty"`
}

func (item *SubpackageRequest) ToModel() *models.Subpackage {
	var availableStartDay *string
	var availableEndDay *string
	if item.IsInf != nil && *item.IsInf {
		availableStartDay = nil
		availableEndDay = nil
	} else {
		availableStartDay = item.AvailableStartDay
		availableEndDay = item.AvailableEndDay
	}
	return &models.Subpackage{
		Title:              *item.Title,
		Description:        *item.Description,
		Price:              *item.Price,
		Duration:           *item.Duration,
		IsInf:              *item.IsInf,
		RepeatedDay:        *item.RepeatedDay,
		AvailableStartTime: *item.AvailableStartTime,
		AvailableEndTime:   *item.AvailableEndTime,
		AvailableStartDay:  utils.SafeString(availableStartDay),
		AvailableEndDay:    utils.SafeString(availableEndDay),
	}
}
