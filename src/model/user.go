package model

import (
	"src/global"
)

// User ...
type User struct {
	Id int64 `json:"id"`
	// Account 账户名
	Username string `json:"username"`
	// Password 密码
	Password string `json:"password"`
	// Time 注册时间
	Time int64 `json:"time"`
	// Status 账户状态：0删除 1正常
	Status int64 `json:"status"`
}

// 通过account获取一个用户
func (self User) OneByAccount() (user User, err error) {
	//global.Db.LogMode(true)
	//log.Println(global.Db.LogMode(true).Error)
	err = global.Dbs.Table("user").Where("username = ?", self.Username).First(&user).Error
	return
}

// 添加用户
func (self User) Add() (err error) {
	err = global.Dbs.Table("user").Create(&self).Error
	return
}

// 通过id获取一个用户
func (self User) One() (user User, err error) {
	err = global.Dbs.Table("user").Where("status = ?", 1).First(&user, self.Id).Error
	return
}
func (self User) One2() (user User, err error) {
	err = global.Dbs.Table("user").First(&user, self.Id).Error
	return
}

// 获取所有用户(不包含管理员、已删除的用户)
func (self User) All() (users []User, err error) {
	err = global.Dbs.Table("user").Where("isAdmin != ?", 1).Where("status != ?", 3).Order("id DESC").Find(&users).Error
	return
}

// 获取所有用户(不包含管理员、已删除的用户)
func (self User) All2() (users []User, err error) {
	err = global.Dbs.Table("user").Where("isAdmin != ?", 1).Order("id DESC").Find(&users).Error
	return
}

// 分页获取用户
func (self User) Page(offset, limit string) (users []User, err error) {
	err = global.Dbs.Table("user").Where("isAdmin != ?", 1).Where("status != ?", 3).Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error
	return
}

// 获取用户总数
func (self User) Total() (total int64, err error) {
	err = global.Dbs.Table("user").Where("status != ?", 3).Count(&total).Error
	return
}

// 修改
func (self User) Edit() (err error) {
	err = global.Dbs.Table("user").Save(&self).Error
	return
}
