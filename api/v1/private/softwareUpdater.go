package private

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kanyuanzhi/middle-platform/global"
	"github.com/kanyuanzhi/middle-platform/model/response"
	pb "github.com/kanyuanzhi/middle-platform/rpc/softwareUpdater"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SoftwareUpdaterApi struct {
	IsUpdating     bool
	LatestVersion  string
	UpdateFilePath []string // 文件路径，包含两个元素，第一个元素是文件所在路径，第二个元素是文件名

	ws *websocket.Conn
}

func (api *SoftwareUpdaterApi) GetSoftwareInfo(c *gin.Context) {
	response.SuccessData(c, global.FXSoftwareInfo)
}

func (api *SoftwareUpdaterApi) CheckUpdateInfo(c *gin.Context) {
	req := &pb.CheckRequest{
		Version:      global.FXSoftwareInfo.Version,
		MachineModel: global.FXSoftwareInfo.MachineModel,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := global.FXSoftwareUpdaterRpcClient.Check(ctxGRPC, req)
	if err != nil {
		log.Printf("gRPC调用失败: %v", err)
		response.ErrorMessage(c, "gRPC调用失败")
		return
	}

	api.LatestVersion = res.GetLatestVersion()
	api.UpdateFilePath = res.GetFilePath()

	checkUpdateInfoResponse := response.CheckUpdateInfo{
		//"isLatest":      res.GetIsLatest(),
		IsLatest:      false,
		LatestVersion: res.GetLatestVersion(),
		HasFile:       res.GetHasFile(),
	}

	response.SuccessData(c, checkUpdateInfoResponse)
}

func (api *SoftwareUpdaterApi) CheckUpdatePermission(c *gin.Context) {
	// 检查控制器是否处在运行状态，运行状态下不允许更新
	var checkUpdatePermissionResponse response.CheckUpdatePermission
	if global.FXControllerStatus.IsRunning || api.IsUpdating {
		checkUpdatePermissionResponse = response.CheckUpdatePermission{
			IsRunning:   global.FXControllerStatus.IsRunning,
			IsUpdating:  api.IsUpdating,
			IsPermitted: false,
		}
	} else {
		checkUpdatePermissionResponse = response.CheckUpdatePermission{
			IsRunning:   global.FXControllerStatus.IsRunning,
			IsUpdating:  api.IsUpdating,
			IsPermitted: true,
		}
	}
	response.SuccessData(c, checkUpdatePermissionResponse)

}

func (api *SoftwareUpdaterApi) Update(c *gin.Context) {
	if api.IsUpdating {
		log.Println("正在更新中，拒绝再次更新")
		response.ErrorMessage(c, errors.New("正在更新中，拒绝再次更新").Error())
		return
	}
	api.IsUpdating = true

	defer func() {
		api.IsUpdating = false
	}()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	log.Println("建立WebSocket连接")
	api.ws = conn
	defer func() {
		conn.Close()
		api.ws = nil
		log.Println("断开WebSocket连接")
	}()

	fileURL := fmt.Sprintf("http://%s:%d/%s", global.FXConfig.SoftwareUpdaterRPC.Host, global.FXConfig.SoftwareUpdaterRPC.FileServerPort,
		strings.Join(api.UpdateFilePath, "/"))
	fmt.Println(fileURL)
	err = api.downloadAndSaveFile(fileURL)
	if err != nil {
		log.Printf("downloadAndSaveFile error:%s", err.Error())
		return
	}

	zipFile := filepath.Join(global.FXConfig.SoftwareUpdaterRPC.SavePath, api.UpdateFilePath[len(api.UpdateFilePath)-1])
	err = api.unzipFile(zipFile)
	if err != nil {
		log.Printf("unzipFile error:%s", err.Error())
		return
	}

	err = api.updateSoftwareInfo()
	if err != nil {
		log.Printf("updateSoftware error:%s", err.Error())
		return
	}
}

func (api *SoftwareUpdaterApi) downloadAndSaveFile(fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// 获取文件的总大小
	totalSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	// 创建本地文件
	file, err := os.Create(filepath.Join(global.FXConfig.SoftwareUpdaterRPC.SavePath, api.UpdateFilePath[len(api.UpdateFilePath)-1]))
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 10240) // 缓冲区大小可以根据需求调整
	startTime := time.Now()
	lastTime := startTime
	lastBytes := 0
	totalBytes := 0

	var downloadSpeed float64 = 0
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err = file.Write(buf[:n])
			if err != nil {
				return err
			}
			totalBytes += n

			currentTime := time.Now()
			elapsedTime := currentTime.Sub(lastTime).Milliseconds()
			if elapsedTime > 1000 {
				downloadSpeed = float64(totalBytes-lastBytes) / (float64(elapsedTime) / 1000) / 1024 / 1024 // MB/s
				lastTime = currentTime
				lastBytes = totalBytes + n
			}

			// 实时发送下载进度到前端
			downloadProgress := float64(totalBytes) / float64(totalSize)
			err = api.sendUpdateData(false, false, downloadProgress, 0, downloadSpeed, 0)
			if err != nil {
				return err
			}

		}
		if err == io.EOF {
			log.Println("下载完毕")
			break
		}
		if err != nil {
			return err
		}
	}
	err = api.sendUpdateData(true, false, 1, 0, 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func (api *SoftwareUpdaterApi) unzipFile(zipFile string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标文件夹
	if err := os.MkdirAll(global.FXConfig.SoftwareUpdaterRPC.UnzipPath, 0755); err != nil {
		return err
	}

	totalFiles := len(r.File)
	completedFiles := 0

	removeUIFolderFlag := false

	// 遍历ZIP文件中的每个文件
	for _, file := range r.File {
		// 构建解压后的文件路径
		extractedFilePath := filepath.Join(global.FXConfig.SoftwareUpdaterRPC.UnzipPath, file.Name)
		// 如果文件是文件夹，创建对应的文件夹
		if file.FileInfo().IsDir() {
			//err = os.MkdirAll(extractedFilePath, file.Mode())
			err = os.MkdirAll(extractedFilePath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if strings.Contains(file.Name, global.FXConfig.SoftwareUpdaterRPC.UIFolderName) && !removeUIFolderFlag {
				removeUIFolderFlag = true
				uiFolderPath := filepath.Join(global.FXConfig.SoftwareUpdaterRPC.SavePath, global.FXConfig.SoftwareUpdaterRPC.UIFolderName)
				log.Printf("发现%s文件夹，删除\n", uiFolderPath)
				err = os.RemoveAll(uiFolderPath)
				if err != nil {
					return err
				}
			}

			if strings.Contains(file.Name, global.FXConfig.SoftwareUpdaterRPC.MiddlePlatformFilename) {
				middlePlatformFilePath := filepath.Join(global.FXConfig.SoftwareUpdaterRPC.SavePath, global.FXConfig.SoftwareUpdaterRPC.MiddlePlatformFilename)
				log.Printf("发现%s文件，删除\n", middlePlatformFilePath)
				err = os.RemoveAll(middlePlatformFilePath)
				if err != nil {
					return err
				}
			}

			if strings.Contains(file.Name, global.FXConfig.SoftwareUpdaterRPC.ControllerFilename) {
				controllerFilePath := filepath.Join(global.FXConfig.SoftwareUpdaterRPC.SavePath, global.FXConfig.SoftwareUpdaterRPC.ControllerFilename)
				log.Printf("发现%s文件，删除\n", controllerFilePath)
				err = os.RemoveAll(controllerFilePath)
				if err != nil {
					return err
				}
			}

			// 创建上层文件夹并解压文件
			if err = os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm); err != nil {
				return err
			}
			// 打开ZIP文件中的文件
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			// 创建目标文件
			dstFile, err := os.Create(extractedFilePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			// 将ZIP文件中的内容复制到目标文件
			_, err = io.Copy(dstFile, rc)
			if err != nil {
				return err
			}

			completedFiles++
			unzipProgress := float64(completedFiles) / float64(totalFiles)
			err = api.sendUpdateData(true, false, 1, unzipProgress, 0, 0)
			if err != nil {
				return err
			}
		}
	}
	log.Println("解压完毕")

	err = api.sendUpdateData(true, true, 1, 1, 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func (api *SoftwareUpdaterApi) sendUpdateData(isDownloadFinished bool, isUnzipFinished bool,
	downloadProgress float64, unzipProgress float64, downloadSpeed float64, unzipSpeed float64) error {
	err := api.ws.WriteJSON(gin.H{
		"isDownloadFinished": isDownloadFinished,
		"isUnzipFinished":    isUnzipFinished,
		"downloadProgress":   downloadProgress,
		"unzipProgress":      unzipProgress,
		"downloadSpeed":      downloadSpeed,
		"unzipSpeed":         unzipSpeed,
	})
	if err != nil {
		log.Printf("error:%s", err.Error())
	}
	return err
}

func (api *SoftwareUpdaterApi) updateSoftwareInfo() error {
	// 读取配置文件
	configFilePath := "softwareInfo.yaml"
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("无法读取配置文件：%v", err)
		return err
	}

	err = yaml.Unmarshal(data, &global.FXSoftwareInfo)
	if err != nil {
		log.Fatalf("无法解析配置文件：%v", err)
		return err
	}

	// 修改字段值
	global.FXSoftwareInfo.Version = api.LatestVersion
	global.FXSoftwareInfo.UpdateTime = time.Now().Local()

	// 将修改后的结构体重新写回配置文件
	newData, err := yaml.Marshal(global.FXSoftwareInfo)
	if err != nil {
		log.Fatalf("无法序列化配置：%v", err)
		return err
	}

	err = os.WriteFile(configFilePath, newData, os.ModePerm)
	if err != nil {
		log.Fatalf("无法写回配置文件：%v", err)
		return err
	}

	log.Println("字段值已修改并写回配置文件")
	return nil
}
