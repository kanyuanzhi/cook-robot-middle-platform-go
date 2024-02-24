package model

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type SysSeasoning struct {
	global.FXModel
	Name   string `json:"name" gorm:"comment:调料名称;"`
	NameEn string `json:"nameEn" gorm:"comment:调料英文名称;"`
	NameTw string `json:"nameTw" gorm:"comment:调料繁体名称;"`
	Pump   uint32 `json:"pump" gorm:"comment:对应泵号"`
	Ratio  uint32 `json:"ratio" gorm:"comment:分量与时间比例"`
}

type SeasoningInfo struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	NameEn string `json:"nameEn"`
	NameTw string `json:"nameTw"`
	Pump   uint32 `json:"pump"`
	Ratio  uint32 `json:"ratio"`
}
