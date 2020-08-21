package model

import (
	"src/global"
)

type Employee struct {
	// Id
	Id int64 `json:"id"`
	// Rid
	Rid int64 `json:"rid"`
	// Name 姓名
	Name string `json:"name"`
	// Sex 性别 1/0
	Sex int64 `json:"sex"`
	// Age 年龄
	Age int64 `json:"age"`
	// Position 职位
	Position string `json:"position"`
	// OnDuty 是否在岗 1/0
	OnDuty int64 `json:"onDuty" gorm:"column:onDuty"`
	// HealthTest 健康检测是否合格 1/0
	HealthTest int64 `json:"healthTest" gorm:"column:healthTest"`
	// TimeTest 健康检测时间
	TimeTest int64 `json:"timeTest" gorm:"column:timeTest"`
	// Temperature 员工每日体温
	Temperature float64 `json:"temperature"`
	// Status 状态
	Status int64 `json:"status"`
}

/**
	分页获取数据
	rid为0时是获取所有的数据，不为0时是获取指定餐厅的员工列表
**/
func (self Employee) Page(rid int64, offset string, limit string) (data []Employee, err error) {
	if rid == 0 {
		err = global.Dbs.Table("employee").Where("status = ?", 1).Order("id DESC").Offset(offset).Limit(limit).Find(&data).Error
	} else {
		err = global.Dbs.Table("employee").Where("status = ?", 1).Where("rid = ?", rid).Order("id DESC").Offset(offset).Limit(limit).Find(&data).Error
	}
	return
}

// 根据餐厅获取数据
func (self Employee) GetDataByRid(rid int64) (data []Employee, err error) {
	if rid == 0 {
		err = global.Dbs.Table("employee").Where("status = ?", 1).Order("id DESC").Find(&data).Error
	} else {
		err = global.Dbs.Table("employee").Where("status = ?", 1).Where("rid = ?", rid).Order("id DESC").Find(&data).Error
	}
	return
}

/**
	获取员工总数
	rid为0时是获取所有的总数，不为0时是获取指定餐厅的员工总数
**/
func (self Employee) Total(rid int64) (total int64, err error) {
	if rid == 0 {
		err = global.Dbs.Table("employee").Where("status = ?", 1).Count(&total).Error
	} else {
		err = global.Dbs.Table("employee").Where("status = ?", 1).Where("rid = ?", rid).Count(&total).Error
	}
	return
}

// 新增一条员工数据
func (self Employee) Add() (err error) {
	err = global.Dbs.Table("employee").Create(&self).Error
	return
}

// 修改数据(包括删除)
func (self Employee) Edit() (err error) {
	err = global.Dbs.Table("employee").Save(&self).Error
	return
}

// 根据id获取一条数据
func (self Employee) One() (user Employee, err error) {
	err = global.Dbs.Table("employee").Where("status = ?", 1).First(&user, self.Id).Error
	return
}
