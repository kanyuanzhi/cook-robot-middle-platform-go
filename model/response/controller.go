package response

import (
	"github.com/kanyuanzhi/middle-platform/global"
)

type FetchControllerStatus struct {
	ControllerStatus global.ControllerStatus `json:"controllerStatus"`
}
