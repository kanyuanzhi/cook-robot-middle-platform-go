package fxcron

import (
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
)

var CronList []model.SysCron
var CronMap = make(map[uuid.UUID]func())
