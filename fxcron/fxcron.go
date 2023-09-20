package fxcron

import (
	"github.com/google/uuid"
	"github.com/kanyuanzhi/middle-platform/model"
)

var CronList []model.SysCron
var CronMap = make(map[uuid.UUID]func())
