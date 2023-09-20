package global

import (
	"github.com/kanyuanzhi/middle-platform/config"
	pb "github.com/kanyuanzhi/middle-platform/rpc/command" // 替换为你的实际包路径
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var (
	FXConfig           config.Config
	FXDb               *gorm.DB
	FXCron             *cron.Cron
	FXCommandRpcClient pb.CommandServiceClient
	FXControllerStatus ControllerStatus
)

type FXModel struct {
	Id        uint   `json:"id" gorm:"comment:id;primaryKey;autoIncrement;"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
	Sort      uint   `json:"sort" gorm:"comment:排序;default:1;"`
	Memo      string `json:"memo" gorm:"comment:备注描述;type:text;"`
}

type InstructionInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Index        int    `json:"index"`
	ActionNumber int    `json:"actionNumber"`
}

type ControllerStatus struct {
	CurrentCommandName              string           `json:"currentCommandName"`
	CurrentDishUUID                 string           `json:"currentDishUUID"`
	CurrentDishCustomStepsUUID      string           `json:"currentDishCustomStepsUUID"`
	CurrentInstructionName          string           `json:"currentInstructionName"`
	CurrentInstructionInfo          *InstructionInfo `json:"currentInstructionInfo"`
	IsPausing                       bool             `json:"isPausing"`
	IsRunning                       bool             `json:"isRunning"`
	IsCooking                       bool             `json:"isCooking"`
	IsPausingWithMovingFinished     bool             `json:"isPausingWithMovingFinished"`
	IsPausingWithMovingBackFinished bool             `json:"isPausingWithMovingBackFinished"`
	IsPausePermitted                bool             `json:"isPausePermitted"`
	BottomTemperature               uint32           `json:"bottomTemperature"`
	InfraredTemperature             uint32           `json:"infraredTemperature"`
	Pump1LiquidWarning              uint32           `json:"pump1LiquidWarning"`
	Pump2LiquidWarning              uint32           `json:"pump2LiquidWarning"`
	Pump3LiquidWarning              uint32           `json:"pump3LiquidWarning"`
	Pump4LiquidWarning              uint32           `json:"pump4LiquidWarning"`
	Pump5LiquidWarning              uint32           `json:"pump5LiquidWarning"`
	Pump6LiquidWarning              uint32           `json:"pump6LiquidWarning"`
	CookingTime                     int64            `json:"cookingTime"`
	CurrentHeatingTemperature       uint32           `json:"currentHeatingTemperature"`
}
