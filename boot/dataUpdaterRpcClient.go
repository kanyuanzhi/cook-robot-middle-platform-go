package boot

import (
	"fmt"
	"github.com/kanyuanzhi/middle-platform/global"
	pb "github.com/kanyuanzhi/middle-platform/rpc/dataUpdater"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DataUpdaterRpcClient() pb.DataUpdaterClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.FXConfig.DataUpdaterRPC.Host, global.FXConfig.DataUpdaterRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(20*1024*1024)))
	if err != nil {
		return nil
	}
	//defer conn.Close()
	client := pb.NewDataUpdaterClient(conn)
	return client
}
