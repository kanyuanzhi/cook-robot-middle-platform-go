package model

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type SysCuisine struct {
	global.FXModel
	Name        string `json:"name" gorm:"comment:菜系名称;"`
	UnDeletable bool   `json:"unDeletable" gorm:"comment:是否允许删除，‘其他’类型不允许删除;default:false"`
}

type CuisineInfo struct {
	Id          uint   `json:"id"`
	Sort        uint   `json:"sort"`
	Name        string `json:"name"`
	UnDeletable bool   `json:"unDeletable"`
}
