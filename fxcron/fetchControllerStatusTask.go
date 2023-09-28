package fxcron

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	pb "github.com/kanyuanzhi/cook-robot-middle-platform-go/rpc/command"
	"log"
	"time"
)

var uid = uuid.New()

var T1 = model.SysCron{
	UUID: uid,
	Id:   0,
	Name: "Fetch Controller Status",
	Spec: "@every 0.1s",
}

func FetchControllerStatus() {
	req := &pb.FetchRequest{
		Empty: true,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := global.FXCommandRpcClient.FetchStatus(ctxGRPC, req)
	if err != nil {
		log.Printf("gRPC调用失败: %v", err)
		return
	}

	var controllerStatus global.ControllerStatus
	err = json.Unmarshal([]byte(res.GetStatusJson()), &controllerStatus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	global.FXControllerStatus = controllerStatus
}
