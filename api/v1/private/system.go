package private

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/response"
	pb "github.com/kanyuanzhi/cook-robot-middle-platform-go/rpc/command"
	"github.com/skip2/go-qrcode"
	"image/png"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type SystemApi struct{}

func (api *SystemApi) GetQrCode(c *gin.Context) {
	ifaces, err := net.Interfaces()
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	// 遍历所有网络接口
	for _, iface := range ifaces {
		// 筛选出WLAN接口，可以根据具体的名称进行判断
		if iface.Name == "wlan0" || iface.Name == "Wi-Fi" || iface.Name == "WLAN" {
			addrs, err := iface.Addrs()
			if err != nil {
				response.ErrorMessage(c, err.Error())
				return
			}

			// 遍历该接口的IP地址
			for _, addr := range addrs {
				// 检查是否是IPv4地址
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					//logger.Log.Println("WLAN IP Address:", ipnet.IP.String())
					qr, err := qrcode.New("phonePairing::"+ipnet.IP.String()+"\r\n", qrcode.Medium)
					if err != nil {
						response.ErrorMessage(c, err.Error())
						return
					}

					var buf bytes.Buffer
					if err := png.Encode(&buf, qr.Image(256)); err != nil {
						response.ErrorMessage(c, err.Error())
						return
					}

					encodedQrImage := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
					getQrCodeResponse := response.GetQrCode{QrCode: encodedQrImage}
					response.SuccessData(c, getQrCodeResponse)
					return
				}
			}
		}
	}
	response.ErrorMessage(c, errors.New("no ip found").Error())
}

func (api *SystemApi) Shutdown(c *gin.Context) {
	req := &pb.ShutdownRequest{
		Empty: true,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, _ := global.FXCommandRpcClient.Shutdown(ctxGRPC, req)
	log.Printf("controller关闭成功%v", res)
	cmd := exec.Command("sudo", "reboot")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
	os.Exit(1)
}
