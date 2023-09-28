package boot

import (
	"fmt"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	pb "github.com/kanyuanzhi/cook-robot-middle-platform-go/rpc/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CommandRpcClient() pb.CommandServiceClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.FXConfig.CommandRPC.Host, global.FXConfig.CommandRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	//defer conn.Close()
	client := pb.NewCommandServiceClient(conn)

	return client
}
