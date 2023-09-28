package utils

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
)

// Reload 解析配置文件configName到target
func Reload(configName string, target interface{}) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("无法读取配置文件:", err)
		return
	}

	err = viper.Unmarshal(target)
	if err != nil {
		log.Println("解析配置文件失败:", err)
		return
	}
}

func LoadLocalImage(path string) ([]byte, error) {
	imagePath := path
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

func GenerateSerialNumber() {
	if global.FXSoftwareInfo.SerialNumber == "" {
		ifaces, err := net.Interfaces()
		if err != nil {
			log.Fatalf("获取网络接口错误：%v", err)
			return
		}

		// 遍历所有网络接口
		for _, iface := range ifaces {
			// 筛选出WLAN接口，可以根据具体的名称进行判断
			if iface.Name == "wlan0" || iface.Name == "Wi-Fi" || iface.Name == "WLAN" {
				global.FXSoftwareInfo.SerialNumber = "SN-XZYC-" + iface.HardwareAddr.String()
				break
			}
		}

		newData, err := yaml.Marshal(global.FXSoftwareInfo)
		if err != nil {
			log.Fatalf("无法序列化配置：%v", err)
		}

		err = os.WriteFile("softwareInfo.yaml", newData, os.ModePerm)
		if err != nil {
			log.Fatalf("无法写回配置文件：%v", err)
		}

		log.Println(global.FXSoftwareInfo)
	} else {
		log.Printf("设备序列号已生成：%s", global.FXSoftwareInfo.SerialNumber)
	}
}
