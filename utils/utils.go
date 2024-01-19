package utils

import (
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm/clause"
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
}

func ResetPersonalDishOwner() {
	var dishes []model.SysDish
	if err := global.FXDb.Where("is_official = ? AND owner != ?", false, global.FXSoftwareInfo.SerialNumber).
		Find(&dishes).Error; err != nil {
		log.Println("ResetPersonalDishOwner失败", err)
		return
	}
	if len(dishes) == 0 {
		log.Println("ResetPersonalDishOwner pass")
		return
	}
	for i := range dishes {
		dishes[i].UUID = uuid.New()
		dishes[i].Owner = global.FXSoftwareInfo.SerialNumber
	}
	if err := global.FXDb.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"uuid", "owner", "updated_at"}),
	}).Create(&dishes).Error; err != nil {
		log.Println("ResetPersonalDishOwner失败", err)
		return
	}
	log.Println("ResetPersonalDishOwner成功")
}
