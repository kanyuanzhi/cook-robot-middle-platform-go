package boot

import (
	"fmt"
	"github.com/kanyuanzhi/middle-platform/fxcron"
	"github.com/kanyuanzhi/middle-platform/global"
	"github.com/robfig/cron/v3"
)

func Cron() {
	c := cron.New(cron.WithSeconds())
	global.FXCron = c
	global.FXCron.Start()
	fmt.Println("Farmoon-Admin Cron Start Succeeded!")
	//global.GqaLogger.Error("Gin-Quasar-Admin Cron Start Succeeded!")
	// Gin-Quasar-Admin cron
	//fxcron.CronList = append(fxcron.CronList, fxcron.T1)
	//fxcron.CronMap[fxcron.T1.UUID] = fxcron.FetchControllerStatus

	for _, task := range fxcron.CronList {
		entryId, _ := global.FXCron.AddFunc(task.Spec, fxcron.CronMap[task.UUID])
		task.Id = entryId
	}
}
