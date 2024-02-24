package model

import (
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type SysDish struct {
	global.FXModel
	Name            string                              `json:"name" gorm:"comment:菜品名称;"`
	UUID            uuid.UUID                           `json:"uuid" gorm:"comment:唯一标识符;not null;uniqueIndex;"`
	Steps           []map[string]interface{}            `json:"steps" gorm:"json;comment:步骤;type:json"`
	CustomStepsList map[string][]map[string]interface{} `json:"customStepsList" gorm:"json;comment:口味步骤;type:json"`
	Image           []byte                              `json:"image" gorm:"comment:菜品图像;"`
	Cuisine         uint                                `json:"cuisine" gorm:"comment:菜系;default:1"`
	IsOfficial      bool                                `json:"isOfficial" gorm:"comment:是否官方菜品;default:false"`
	IsShared        bool                                `json:"isShared" gorm:"comment:是否共享菜品;default:false"`
	IsMarked        bool                                `json:"isMarked" gorm:"是否收藏菜品至我的菜品;default:false"`
	Owner           string                              `json:"owner" gorm:"comment:菜品所有者;default:''"`
	Local           string                              `json:"local" gorm:"comment:本地语言（en,cn,tw）;default:'cn'"`
}

type DishInfo struct {
	Id              uint                                `json:"id"`
	Name            string                              `json:"name"`
	UUID            uuid.UUID                           `json:"uuid"`
	Steps           []map[string]interface{}            `json:"steps"`
	CustomStepsList map[string][]map[string]interface{} `json:"customStepsList"`
	Image           string                              `json:"image"`
	Cuisine         uint                                `json:"cuisine"`
	IsOfficial      bool                                `json:"isOfficial"`
	IsShared        bool                                `json:"isShared"`
	IsMarked        bool                                `json:"isMarked"`
	Owner           string                              `json:"owner"`
	Local           string                              `json:"local"`
}
