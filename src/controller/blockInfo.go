package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"src/model"
)

type BlockInfo struct {
	IsNew bool
	Info  []Infos
}
type Infos struct {
	Block  Blocks
	Status []Statuses
}
type Blocks struct {
	Height int64
	Hash   string
	Time   string
}
type Statuses struct {
	Deviename string
	Des       string
	Time      string
}

func BlockView(c *gin.Context) {
	var respData RespData
	var data BlockInfo
	data.IsNew = true
	// 接收数据
	heightStr := c.Query("height")
	// 验证数据
	if heightStr == "" {
		respData.Error = "区块高度不能为空"
		c.JSON(http.StatusOK, &respData)
		return
	}
	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		respData.Error = "区块高度必须是数字"
		c.JSON(http.StatusOK, &respData)
		return
	}
	fmt.Println("height: ", height)

	// 获取最新区块高度
	var block model.BlockInfo
	maxHeight, err := block.MaxHeight()
	if err != nil {
		respData.Error = "最新区块高度不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}

	if height != 0 {
		// 验证当前高度是否为最新区块高度
		block.Height = height
		blockOld, err := block.GetDataByHeight()
		if err != nil {
			log.Println(err)
			respData.Error = "区块高度不存在"
			c.JSON(http.StatusOK, &respData)
			return
		}
		// 判断是否最新区块
		if blockOld.Height < maxHeight {
			data.IsNew = true
		} else {
			data.IsNew = false
		}
	}

	// 获取所有设备名称
	var device model.DeviceInfo
	deviceInfo, _ := device.GetAll()
	deviceMap := make(map[string]string)
	for _, v := range deviceInfo {
		deviceMap[fmt.Sprint(v.Parentid)+"_"+fmt.Sprint(v.Channel)] = v.Devicename
	}
	// 获取所有状态类型说明
	var st model.StatusType
	statusTypes, _ := st.All()
	statusTypesMap := make(map[int64]string)
	for _, v := range statusTypes {
		statusTypesMap[v.Id] = v.Des
	}

	// 通过区块高度获取数据
	var infos []Infos
	for i := 0; i < 3; i++ {
		// 第一个区块数据
		block.Height = maxHeight - int64(i)
		currentBlock, _ := block.GetDataByHeight()

		// 第二个区块数据
		block.Height = maxHeight - int64(i) - 1
		previousBlock, _ := block.GetDataByHeight()

		var info Infos
		var a_block Blocks
		a_block.Height = currentBlock.Height
		a_block.Hash = currentBlock.Hash
		a_block.Time = currentBlock.Time
		info.Block = a_block

		if len([]rune(currentBlock.Time)) >= 10 && len([]rune(previousBlock.Time)) >= 10 {
			// 通过时间范围获取设备状态
			var ds model.DeviceStatus
			statuses, err := ds.GetDataByTime(previousBlock.Time, currentBlock.Time)
			if err != nil {
				respData.Error = "获取数据失败"
				c.JSON(http.StatusOK, &respData)
				return
			}
			var statu []Statuses
			for _, v := range statuses {
				var s Statuses
				s.Deviename = deviceMap[fmt.Sprint(v.DeviceId)+"_"+fmt.Sprint(v.Channel)]
				s.Des = statusTypesMap[v.StatusTypeId]
				s.Time = v.Time
				statu = append(statu, s)
			}
			info.Status = statu
		} else {
			info.Status = nil
		}
		infos = append(infos, info)
	}
	data.Info = infos
	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}
