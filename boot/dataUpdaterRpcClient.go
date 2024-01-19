package boot

import (
	"fmt"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	pb "github.com/kanyuanzhi/cook-robot-middle-platform-go/rpc/dataUpdater"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DataUpdaterRpcClient() pb.DataUpdaterClient {
	maxSize := 100 * 1024 * 1024

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.FXConfig.DataUpdaterRPC.Host, global.FXConfig.DataUpdaterRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxSize), grpc.MaxCallSendMsgSize(maxSize)))
	if err != nil {
		return nil
	}
	//defer conn.Close()
	client := pb.NewDataUpdaterClient(conn)
	return client
}
