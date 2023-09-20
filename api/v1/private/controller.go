package private

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/middle-platform/controller/command"
	"github.com/kanyuanzhi/middle-platform/controller/instruction"
	"github.com/kanyuanzhi/middle-platform/global"
	"github.com/kanyuanzhi/middle-platform/model"
	"github.com/kanyuanzhi/middle-platform/model/request"
	"github.com/kanyuanzhi/middle-platform/model/response"
	pb "github.com/kanyuanzhi/middle-platform/rpc/command"
	"github.com/mitchellh/mapstructure"
	"time"
)

type ControllerApi struct{}

func (api *ControllerApi) Execute(c *gin.Context) {
	var executeCommandRequest request.ExecuteCommandRequest

	if err := request.ShouldBindJSON(c, &executeCommandRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var commandStruct command.Command
	if executeCommandRequest.CommandType == command.COMMAND_TYPE_MULTIPLE {
		// 多指令
		if executeCommandRequest.CommandName == command.COMMAND_NAME_COOK {
			var dish model.SysDish
			err := global.FXDb.First(&dish, "uuid = ?", executeCommandRequest.CommandData.UUID).Error
			if err != nil {
				response.ErrorMessage(c, err.Error())
				return
			}
			var steps []map[string]interface{}
			if executeCommandRequest.CommandData.CustomStepsUUID != "" {
				steps = dish.CustomStepsList[executeCommandRequest.CommandData.CustomStepsUUID]
			} else {
				steps = dish.Steps
			}
			//logger.Log.Println(dbDish)
			//var stepsJSON []map[string]interface{}
			//err = json.Unmarshal([]byte(dish.Steps), &stepsJSON)
			//if err != nil {
			//	response.ErrorMessage(c, err.Error())
			//	return
			//}

			var seasonings []model.SysSeasoning
			err = global.FXDb.Select("pump", "ratio").Find(&seasonings).Error
			if err != nil {
				response.ErrorMessage(c, err.Error())
				return
			}
			pumpToRatioMap := map[string]uint32{}
			for _, seasoning := range seasonings {
				pumpToRatioMap[fmt.Sprintf("%d", seasoning.Pump)] = seasoning.Ratio
			}

			var instructions []instruction.Instructioner
			// 开始先启动转动、油烟净化
			instructions = append(instructions, instruction.NewInitInstruction("启动中"))

			for _, step := range steps {
				instructionType := instruction.InstructionType(step["instructionType"].(string))
				var instructionStruct instruction.Instructioner
				if instructionType == instruction.SEASONING {
					pumpToWeightMap := map[string]uint32{}
					for _, seasoning := range step["seasonings"].([]interface{}) {
						pumpNumber := fmt.Sprintf("%.0f", seasoning.(map[string]interface{})["pumpNumber"].(float64))
						pumpToWeightMap[pumpNumber] = uint32(seasoning.(map[string]interface{})["weight"].(float64))
					}
					instructionStruct = instruction.NewSeasoningInstruction(step["instructionName"].(string), pumpToWeightMap, pumpToRatioMap)
				} else {
					instructionStruct = instruction.InstructionTypeToStruct[instructionType]
					err := mapstructure.Decode(step, &instructionStruct)
					if err != nil {
						response.ErrorMessage(c, err.Error())
						return
					}
					if instructionType == instruction.WATER {
						if t, ok := instructionStruct.(instruction.WaterInstruction); ok {
							t.Ratio = pumpToRatioMap["6"]
							instructions = append(instructions, t)
							continue
						}
					}
					if instructionType == instruction.OIL {
						if t, ok := instructionStruct.(instruction.OilInstruction); ok {
							t.Ratio = pumpToRatioMap["1"]
							instructions = append(instructions, t)
							continue
						}
					}
				}
				instructions = append(instructions, instructionStruct)
			}

			instructions = append(instructions, instruction.NewFinishInstruction("停止中"))

			commandStruct = command.Command{
				CommandName:     command.COMMAND_NAME_COOK,
				CommandType:     command.COMMAND_TYPE_MULTIPLE,
				DishUUID:        executeCommandRequest.CommandData.UUID,
				CustomStepsUUID: executeCommandRequest.CommandData.CustomStepsUUID,
				Instructions:    instructions,
			}

		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_PREPARE {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewPrepareInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_PREPARE,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_DISH_OUT {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewDishOutInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_DISH_OUT,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_WASH {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewWashInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_WASH,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_POUR {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewPourInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_POUR,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_WITHDRAW {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewWithdrawInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_WITHDRAW,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		} else {
			response.ErrorMessage(c, errors.New(fmt.Sprintf("%s指令错误", executeCommandRequest.CommandName)).Error())
			return
		}
	} else {
		// 单指令，立即执行
		if executeCommandRequest.CommandName == command.COMMAND_NAME_DOOR_UNLOCK {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewDoorUnlockInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_DOOR_UNLOCK,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_PAUSE_TO_ADD {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewPauseToAddInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_PAUSE_TO_ADD,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_RESUME {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewResumeInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_RESUME,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else if executeCommandRequest.CommandName == command.COMMAND_NAME_HEAT {
			var instructions []instruction.Instructioner
			//temperature, err := strconv.ParseFloat(executeCommandRequest.CommandData, 10)
			//if err != nil {
			//	response.ErrorMessage(c, err.Error())
			//	return
			//}
			instructions = append(instructions, instruction.NewHeatInstruction(
				executeCommandRequest.CommandData.Temperature, 0, 0, instruction.NO_JUDGE))
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_HEAT,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else {
			response.ErrorMessage(c, errors.New("命令名称错误").Error())
			return
		}
	}

	commandJSON, err := json.Marshal(commandStruct)

	req := &pb.CommandRequest{
		CommandJson: string(commandJSON),
	}

	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := global.FXCommandRpcClient.Execute(ctxGRPC, req)

	if err != nil {
		response.ErrorMessage(c, errors.New(fmt.Sprintf("gRPC调用失败: %v", err)).Error())
		return
	}

	if res.GetResult() == 0 {
		response.ErrorMessage(c, errors.New("机器占用中").Error())
		return
	}

	response.SuccessMessage(c, "执行成功")
}

func (api *ControllerApi) FetchStatus(c *gin.Context) {
	fetchControllerStatusResponse := response.FetchControllerStatus{
		ControllerStatus: global.FXControllerStatus,
	}
	response.SuccessData(c, fetchControllerStatusResponse)
}
