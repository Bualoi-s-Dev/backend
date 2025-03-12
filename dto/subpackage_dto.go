package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/utils"
)

type SubpackageRequest struct {
	Title       *string `bson:"title" json:"title" binding:"omitempty" example:"Wedding Bliss Package"`
	Description *string `bson:"description" json:"description" binding:"omitempty" example:"This is a package for wedding"`
	Price       *int    `bson:"price" json:"price" binding:"omitempty" example:"10000"`
	Duration    *int    `bson:"duration" json:"duration" binding:"omitempty" example:"60" description:"Duration in minutes"`

	IsInf              *bool             `bson:"is_inf" json:"isInf" binding:"omitempty,isInf_rule" example:"false"`
	RepeatedDay        *[]models.DayName `bson:"repeated_day" json:"repeatedDay" binding:"omitempty,day_names" example:"MON,TUE,WED"`
	AvaliableStartTime *string           `bson:"avaliable_start_time" json:"avaliableStartTime" binding:"omitempty,time_format" example:"15:04"`
	AvaliableEndTime   *string           `bson:"avaliable_end_time" json:"avaliableEndTime" binding:"omitempty,time_format" example:"16:27"`
	AvaliableStartDay  *string           `bson:"avaliable_start_day" json:"avaliableStartDay" binding:"omitempty,date_format" example:"2021-01-01"`
	AvaliableEndDay    *string           `bson:"avaliable_end_day" json:"avaliableEndDay" binding:"omitempty,date_format" example:"2021-12-31"`
}

func (item *SubpackageRequest) ToModel() *models.Subpackage {
	var avaliableStartDay *string
	var avaliableEndDay *string
	if item.IsInf != nil && *item.IsInf {
		avaliableStartDay = nil
		avaliableEndDay = nil
	} else {
		avaliableStartDay = item.AvaliableStartDay
		avaliableEndDay = item.AvaliableEndDay
	}
	return &models.Subpackage{
		Title:              *item.Title,
		Description:        *item.Description,
		Price:              *item.Price,
		Duration:           *item.Duration,
		IsInf:              *item.IsInf,
		RepeatedDay:        *item.RepeatedDay,
		AvaliableStartTime: *item.AvaliableStartTime,
		AvaliableEndTime:   *item.AvaliableEndTime,
		AvaliableStartDay:  utils.SafeString(avaliableStartDay),
		AvaliableEndDay:    utils.SafeString(avaliableEndDay),
	}
}
