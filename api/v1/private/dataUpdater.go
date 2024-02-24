package private

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/response"
	pb "github.com/kanyuanzhi/cook-robot-middle-platform-go/rpc/dataUpdater"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"time"
)

type DataUpdaterApi struct{}

func (api *DataUpdaterApi) UpdateOfficialDishes(c *gin.Context) {
	cmd := exec.Command("sudo", "nmcli device modify eth0 ipv4.route-metric 1000")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
	}

	if global.FXControllerStatus.IsRunning {
		response.ErrorMessage(c, "有菜品正在炒制，请稍后同步菜品")
		return
	}

	var officialDishes []model.SysDish
	if err := global.FXDb.Where("is_official = ?", true).Find(&officialDishes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}
	officialDishesInfo := make(map[uuid.UUID]int64)
	for _, officialDish := range officialDishes {
		officialDishesInfo[officialDish.UUID] = officialDish.UpdatedAt
	}
	officialDishesInfoBytes, _ := json.Marshal(officialDishesInfo)

	req := &pb.FetchOfficialDishesRequest{
		LocalDishesInfoJson: officialDishesInfoBytes,
		UserSerialNumber:    global.FXSoftwareInfo.SerialNumber,
		Version:             global.FXSoftwareInfo.Version,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := global.FXDataUpdaterRpcClient.FetchOfficialDishes(ctxGRPC, req)
	if err != nil {
		log.Println(err)
		response.ErrorMessage(c, "RPC调用失败【fetchOfficials】，"+err.Error())
		return
	}

	var needAddDishes []model.SysDish
	var needUpdateDishes []model.SysDish
	var needDeleteDishUUIDs []uuid.UUID
	var cuisines []model.SysCuisine
	if err := json.Unmarshal(res.GetLocalNeedAddDishesJson(), &needAddDishes); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}
	if err := json.Unmarshal(res.GetLocalNeedUpdateDishesJson(), &needUpdateDishes); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}
	if err := json.Unmarshal(res.GetLocalNeedDeleteDishesUuidsJson(), &needDeleteDishUUIDs); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}
	if err := json.Unmarshal(res.GetCuisinesJson(), &cuisines); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}

	tx := global.FXDb.Begin()
	if len(needAddDishes) != 0 {
		if err := tx.Create(needAddDishes).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if len(needUpdateDishes) != 0 {
		for _, dish := range needUpdateDishes {
			if err := tx.Model(&model.SysDish{}).Where("uuid = ?", dish.UUID).Omit("id", "uuid").Updates(dish).Error; err != nil {
				tx.Rollback()
				response.ErrorMessage(c, err.Error())
				return
			}
		}
	}

	if len(needDeleteDishUUIDs) != 0 {
		if err := tx.Where("uuid in ?", needDeleteDishUUIDs).Delete(&model.SysDish{}).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if len(cuisines) != 0 {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.SysCuisine{}).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
		if err := tx.Create(cuisines).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	api.UpdateIngredients(c)

	updatedOfficialDishesResponse := response.UpdateOfficialDishes{
		NewAddedDishesNumber: len(needAddDishes),
		UpdatesDishesNumber:  len(needUpdateDishes),
		DeletedDishesNumber:  len(needDeleteDishUUIDs),
	}

	response.SuccessMessageData(c, updatedOfficialDishesResponse, "更新成功")
}

func (api *DataUpdaterApi) UpdateIngredients(c *gin.Context) {
	req := &pb.FetchIngredientsRequest{
		UserSerialNumber: global.FXSoftwareInfo.SerialNumber,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := global.FXDataUpdaterRpcClient.FetchIngredients(ctxGRPC, req)
	if err != nil {
		response.ErrorMessage(c, "RPC调用失败")
		return
	}

	var ingredients []model.SysIngredient
	var ingredientTypes []model.SysIngredientType
	var ingredientShapes []model.SysIngredientShape
	if err := json.Unmarshal(res.GetIngredientsJson(), &ingredients); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}
	if err := json.Unmarshal(res.GetIngredientTypesJson(), &ingredientTypes); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}
	if err := json.Unmarshal(res.GetIngredientShapesJson(), &ingredientShapes); err != nil {
		log.Printf("json unmarshal失败: %v", err)
	}

	tx := global.FXDb.Begin()

	if len(ingredients) != 0 {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.SysIngredient{}).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
		if err := tx.Create(ingredients).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if len(ingredientTypes) != 0 {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.SysIngredientType{}).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
		if err := tx.Create(ingredientTypes).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if len(ingredientShapes) != 0 {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.SysIngredientShape{}).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
		if err := tx.Create(ingredientShapes).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	//response.SuccessMessage(c, "更新成功")
}

func (api *DataUpdaterApi) SynchronizePersonalDishes(c *gin.Context) {
	var personalDishes []model.SysDish
	if err := global.FXDb.Where("is_official = ? AND owner = ?", false, global.FXSoftwareInfo.SerialNumber).Find(&personalDishes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}
	personalDishesInfo := make(map[uuid.UUID]int64)
	for _, personalDish := range personalDishes {
		personalDishesInfo[personalDish.UUID] = personalDish.UpdatedAt
	}
	personalDishesInfoBytes, _ := json.Marshal(personalDishesInfo)

	var userDeletedDishes []model.SysUserDeletedDish
	if err := global.FXDb.Find(&userDeletedDishes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}
	var localDeletedDishUUIDs []uuid.UUID
	for _, userDeletedDish := range userDeletedDishes {
		localDeletedDishUUIDs = append(localDeletedDishUUIDs, userDeletedDish.UUID)
	}
	localDeletedDishUuidsJson, _ := json.Marshal(localDeletedDishUUIDs)

	synchronizeReq := &pb.SynchronizePersonalDishesRequest{
		UserSerialNumber:          global.FXSoftwareInfo.SerialNumber,
		LocalDishesInfoJson:       personalDishesInfoBytes,
		LocalDeletedDishUuidsJson: localDeletedDishUuidsJson,
		Version:                   global.FXSoftwareInfo.Version,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	synchronizeRes, err := global.FXDataUpdaterRpcClient.SynchronizePersonalDishes(ctxGRPC, synchronizeReq)
	if err != nil {
		response.ErrorMessage(c, "RPC调用失败【synchronize】，"+err.Error())
		return
	}

	// 远程已经删除用户在本地删除的菜品，需要将本地删除的菜品从数据库中删除
	if err := global.FXDb.Where("uuid in ?", localDeletedDishUUIDs).Delete(&model.SysUserDeletedDish{}).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var remoteNeedAddDishUUIDs []uuid.UUID
	var remoteNeedUpdateDishUUIDs []uuid.UUID
	var localNeedAddDishes []model.SysDish
	var localNeedUpdateDishes []model.SysDish
	var localNeedDeleteDishUUIDs []uuid.UUID

	_ = json.Unmarshal(synchronizeRes.GetRemoteNeedAddDishUuidsJson(), &remoteNeedAddDishUUIDs)
	_ = json.Unmarshal(synchronizeRes.GetRemoteNeedUpdateDishUuidsJson(), &remoteNeedUpdateDishUUIDs)
	_ = json.Unmarshal(synchronizeRes.GetLocalNeedAddDishesJson(), &localNeedAddDishes)
	_ = json.Unmarshal(synchronizeRes.GetLocalNeedUpdateDishesJson(), &localNeedUpdateDishes)
	_ = json.Unmarshal(synchronizeRes.GetLocalNeedDeleteDishUuidsJson(), &localNeedDeleteDishUUIDs)

	var remoteNeedAddDishes []model.SysDish
	var remoteNeedUpdateDishes []model.SysDish
	tx := global.FXDb.Begin()
	if len(remoteNeedAddDishUUIDs) != 0 {
		if err := tx.Where("uuid in ?", remoteNeedAddDishUUIDs).Omit("id").Find(&remoteNeedAddDishes).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}
	if len(remoteNeedUpdateDishUUIDs) != 0 {
		if err := tx.Where("uuid in ?", remoteNeedUpdateDishUUIDs).Omit("id").Find(&remoteNeedUpdateDishes).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if len(localNeedAddDishes) != 0 {
		if err := tx.Create(localNeedAddDishes).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if len(localNeedUpdateDishes) != 0 {
		for _, dish := range localNeedUpdateDishes {
			if err := tx.Model(&model.SysDish{}).Where("uuid = ?", dish.UUID).Omit("id", "uuid").Updates(dish).Error; err != nil {
				tx.Rollback()
				response.ErrorMessage(c, err.Error())
				return
			}
			// 将自动更新的updated_at字段更新为远端下载过来的值
			if err := tx.Model(&model.SysDish{}).Where("uuid = ?", dish.UUID).Update("updated_at", dish.UpdatedAt).Error; err != nil {
				tx.Rollback()
				response.ErrorMessage(c, err.Error())
				return
			}
		}
	}

	if len(localNeedDeleteDishUUIDs) != 0 {
		if err := tx.Where("uuid in ?", localNeedDeleteDishUUIDs).Delete(&model.SysDish{}).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	remoteNeedAddDishesBytes, _ := json.Marshal(remoteNeedAddDishes)
	remoteNeedUpdateDishesBytes, _ := json.Marshal(remoteNeedUpdateDishes)

	uploadReq := &pb.UploadPersonalDishesRequest{
		UserSerialNumber:           global.FXSoftwareInfo.SerialNumber,
		RemoteNeedAddDishesJson:    remoteNeedAddDishesBytes,
		RemoteNeedUpdateDishesJson: remoteNeedUpdateDishesBytes,
	}
	_, err = global.FXDataUpdaterRpcClient.UploadPersonalDishes(ctxGRPC, uploadReq)
	if err != nil {
		response.ErrorMessage(c, "RPC调用失败【upload】"+err.Error())
		return
	}

	synchronizePersonalDishesResponse := response.SynchronizePersonalDishes{
		RemoteNeedAddDishesNumber:    len(remoteNeedAddDishUUIDs),
		RemoteNeedUpdateDishesNumber: len(remoteNeedUpdateDishUUIDs),
		RemoteNeedDeleteDishesNumber: int(synchronizeRes.GetRemoteNeedDeleteDishesNumber()),
		LocalNeedAddDishesNumber:     len(localNeedAddDishes),
		LocalNeedUpdateDishesNumber:  len(localNeedUpdateDishes),
		LocalNeedDeleteDishesNumber:  len(localNeedDeleteDishUUIDs),
	}

	response.SuccessMessageData(c, synchronizePersonalDishesResponse, "同步成功")
}
