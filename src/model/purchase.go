package model

import (
	"errors"
	"src/global"
)

type Purchase struct {
	// Id
	Id int64 `json:"id"`
	// Rid
	Rid int64 `json:"rid"`
	// Time 记录时间
	Time int64 `json:"time"`
	// Goods 物品数据
	Goods string `json:"goods"`
	// Status 状态
	Status int64 `json:"status"`
}

/**
	分页获取指定餐厅内采购物品数据
	rid 必须大于0
	start/end 必须同时都在 且end>=start 如果不带这两个参数 就都传0
**/
func (self Purchase) Page(rid int64, start int64, end int64, offset string, limit string) (data []Purchase, err error) {
	if rid > 0 {
		if start == 0 && end == 0 {
			err = global.Dbs.Table("purchase").Where("status = ?", 1).Where("rid = ?", rid).Order("id DESC").Offset(offset).Limit(limit).Find(&data).Error
		} else {
			if start > 0 && end >= start {
				err = global.Dbs.Table("purchase").Where("status = ?", 1).Where("rid = ?", rid).Where("time BETWEEN ? AND ?", start, end).Order("id DESC").Offset(offset).Limit(limit).Find(&data).Error
			} else {
				err = errors.New("提交的时间查询参数错误")
			}
		}
	} else {
		err = errors.New("提交的餐厅参数错误")
	}
	return
}

// 根据餐厅获取数据
func (self Purchase) GetDataByRid(rid int64, start int64, end int64) (data []Purchase, err error) {
	if rid > 0 {
		if start == 0 && end == 0 {
			err = global.Dbs.Table("purchase").Where("status = ?", 1).Where("rid = ?", rid).Order("id DESC").Find(&data).Error
		} else {
			if start > 0 && end >= start {
				err = global.Dbs.Table("purchase").Where("status = ?", 1).Where("rid = ?", rid).Where("time BETWEEN ? AND ?", start, end).Order("id DESC").Find(&data).Error
			} else {
				err = errors.New("提交的时间查询参数错误")
			}
		}
	} else {
		err = errors.New("提交的餐厅参数错误")
	}
	return
}

/**
	获取指定餐厅内物品采购总数
	rid必须大于0
	start/end 必须同时都在 且end>=start 如果不带这两个参数 就都传0
**/
func (self Purchase) Total(rid int64, start int64, end int64) (total int64, err error) {
	total = 0
	if rid > 0 {
		if start == 0 && end == 0 {
			err = global.Dbs.Table("purchase").Where("status = ?", 1).Where("rid = ?", rid).Count(&total).Error
		} else {
			if start > 0 && end >= start {
				err = global.Dbs.Table("purchase").Where("status = ?", 1).Where("rid = ?", rid).Where("time BETWEEN ? AND ?", start, end).Count(&total).Error
			} else {
				err = errors.New("提交的时间查询参数错误")
			}
		}
	} else {
		err = errors.New("提交的餐厅参数错误")
	}
	return
}

// 新增一条采购数据
func (self Purchase) Add() (err error) {
	err = global.Dbs.Table("purchase").Create(&self).Error
	return
}

// 修改数据(包括删除)
func (self Purchase) Edit() (err error) {
	err = global.Dbs.Table("purchase").Save(&self).Error
	return
}

// 根据id获取一条数据
func (self Purchase) One() (user Purchase, err error) {
	err = global.Dbs.Table("purchase").Where("status = ?", 1).First(&user, self.Id).Error
	return
}
