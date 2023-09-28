package response

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type FetchControllerStatus struct {
	ControllerStatus global.ControllerStatus `json:"controllerStatus"`
}
