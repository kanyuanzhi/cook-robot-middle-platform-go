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
	"strings"
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

func AddColumns() {
	addColumn("sys_dish", "local")

	updateLocalSql := "UPDATE sys_dish SET local = 'cn' WHERE local is NULL;"
	if err := global.FXDb.Exec(updateLocalSql).Error; err != nil {
		log.Println("AddColumn失败", err)
		return
	}

	addColumn("sys_cuisine", "name_en")
	addColumn("sys_cuisine", "name_tw")
	addColumn("sys_ingredient", "name_en")
	addColumn("sys_ingredient", "name_tw")
	addColumn("sys_ingredient_type", "name_en")
	addColumn("sys_ingredient_type", "name_tw")
	addColumn("sys_ingredient_shape", "name_en")
	addColumn("sys_ingredient_shape", "name_tw")
	addColumn("sys_seasoning", "name_en")
	addColumn("sys_seasoning", "name_tw")

	sql := "UPDATE sys_seasoning " +
		"SET name_en = " +
		"CASE " +
		"WHEN name = '食用油' THEN 'oil' " +
		"WHEN name = '生抽' THEN 'light soy sauce' " +
		"WHEN name = '老抽' THEN 'soy sauce' " +
		"WHEN name = '醋' THEN 'vinegar' " +
		"WHEN name = '料酒' THEN 'cooking wine' " +
		"WHEN name = '纯净水' THEN 'water' " +
		"WHEN name = '自来水1' THEN 'tap water1' " +
		"WHEN name = '自来水2' THEN 'tap water2' " +
		"WHEN name = '食盐' THEN 'salt' " +
		"WHEN name = '鸡精' THEN 'chicken powder' " +
		"ELSE name_en " +
		"END, " +
		"name_tw = " +
		"CASE " +
		"WHEN name = '食用油' THEN '食用油' " +
		"WHEN name = '生抽' THEN '生抽' " +
		"WHEN name = '老抽' THEN '老抽' " +
		"WHEN name = '醋' THEN '醋' " +
		"WHEN name = '料酒' THEN '料酒' " +
		"WHEN name = '纯净水' THEN '純淨水' " +
		"WHEN name = '自来水1' THEN '自來水1' " +
		"WHEN name = '自来水2' THEN '自來水2' " +
		"WHEN name = '食盐' THEN '食鹽' " +
		"WHEN name = '鸡精' THEN '鸡精' " +
		"ELSE name_tw " +
		"END;"
	if err := global.FXDb.Exec(sql).Error; err != nil {
		log.Println("AddColumn失败", err)
		return
	}
	log.Println("AddColumns成功")
}

func addColumn(table string, column string) error {
	sql := "ALTER TABLE " + table + " ADD " + column + " TEXT;"
	if err := global.FXDb.Exec(sql).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate column name") {
			log.Println(column + "字段已存在")
		} else {
			log.Println("AddColumn失败", err)
			return err
		}
	}
	return nil
}
