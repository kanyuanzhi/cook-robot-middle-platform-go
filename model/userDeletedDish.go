package model

import (
	"github.com/google/uuid"
	"github.com/kanyuanzhi/middle-platform/global"
)

type SysUserDeletedDish struct {
	global.FXModel
	UUID  uuid.UUID `json:"uuid" gorm:"comment:唯一标识符;not null;uniqueIndex;"`
	Owner string    `json:"owner" gorm:"comment:菜品所有者;default:''"`
}
