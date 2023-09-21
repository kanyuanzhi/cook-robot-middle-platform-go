package boot

import (
	"fmt"
	"github.com/kanyuanzhi/middle-platform/global"
	pb "github.com/kanyuanzhi/middle-platform/rpc/softwareUpdater"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SoftwareUpdaterRpcClient() pb.UpdateClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.FXConfig.SoftwareUpdaterRPC.Host, global.FXConfig.SoftwareUpdaterRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	//defer conn.Close()
	client := pb.NewUpdateClient(conn)
	return client
}
