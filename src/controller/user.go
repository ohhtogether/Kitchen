package controller

import (
	"log"
	"net/http"
	"strconv"
	//"strconv"

	"github.com/dxvgef/filter"
	"github.com/gin-gonic/gin"

	"src/global"
	"src/model"
)

/*// 注册
func Register(c *gin.Context) {
	var respData RespData

	// 接受数据
	account := c.PostForm("account")
	password := c.PostForm("password")
	confirmPwd := c.PostForm("confirmPwd")
	name := c.PostForm("name")
	genderStr := c.PostForm("gender")
	phoneStr := c.PostForm("phone")
	identityIdStr := c.PostForm("identityId")

	// 验证
	err := filter.MSet(
		filter.El(&account,
			filter.FromString(account, "账户名").
				Required().MinLength(6).MaxLength(32).IsLetterOrDigit(),
		),
	)
	if err != nil {
		respData.Code = 1011
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = filter.MSet(
		filter.El(&password,
			filter.FromString(password, "登录密码").
				Required().MinLength(6).MaxLength(32).IsLetterOrDigit(),
		),
	)
	if err != nil {
		respData.Code = 1012
		c.JSON(http.StatusOK, &respData)
		return
	}
	if password != confirmPwd {
		respData.Code = 1023
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = filter.MSet(
		filter.El(&name,
			filter.FromString(name, "姓名").
				Required().MinLength(4).MaxLength(16).IsChinese(),
		),
	)
	if err != nil {
		respData.Code = 1013
		c.JSON(http.StatusOK, &respData)
		return
	}
	gender, err := strconv.ParseInt(genderStr, 10, 64)
	if err != nil || (gender != 0 && gender != 1 && gender != 2) {
		respData.Code = 1014
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = global.VerifyMobileFormat(phoneStr)
	if err != nil {
		respData.Code = 1015
		c.JSON(http.StatusOK, &respData)
		return
	}
	phone, err := strconv.ParseInt(phoneStr, 10, 64)
	if err != nil {
		respData.Code = 1015
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = global.IsValidCitizenNo18(identityIdStr)
	if err != nil {
		respData.Code = 1016
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 验证账号是否存在
	var user model.User
	user.Account = account
	userOne, err := user.OneByAccount()
	if err == nil && userOne.Account == account {
		respData.Code = 1010
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 密码加密
	hashpass, err := global.HashPassword(password)
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	user.Name = name
	user.Gender = gender
	user.Phone = phone
	user.IdentityId = identityIdStr
	user.Time = time.Now().Unix()
	user.Account = account
	user.Password = hashpass
	err = user.Add()
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	c.JSON(http.StatusOK, &respData)
	return
}*/

// 登录
func Login(c *gin.Context) {
	var respData RespData

	// 接受数据
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证
	err := filter.MSet(
		filter.El(&username,
			filter.FromString(username, "登录账号").
				Required().MinLength(5).MaxLength(32).IsLetterOrDigit(),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}
	if global.FilteredSQLInject(username) {
		respData.Error = "账号含非法字符"
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = filter.MSet(
		filter.El(&password,
			filter.FromString(password, "登录密码").
				Required().MinLength(5).MaxLength(32),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}
	//hashpass, err := global.HashPassword(password)

	var user model.User
	user.Username = username
	userOne, err := user.OneByAccount()
	if err != nil {
		log.Println(err)
		respData.Error = "用户不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 验证密码
	Bool := global.CheckPasswordHash(password, userOne.Password)
	if !Bool {
		respData.Error = "密码错误!"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 生成Token
	var userT global.UserToken
	userT.Id = userOne.Id
	userT.Name = userOne.Username
	token, err := global.GenToken(userT)
	if err != nil {
		log.Println(err)
		respData.Error = "登录失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	result := struct {
		Id      int64  `json:"id"`
		Name    string `json:"name"`
		Token   string `json:"token"`
		Isadmin int64  `json:"isadmin"`
	}{
		Id:      userOne.Id,
		Name:    userOne.Username,
		Token:   token,
	}

	respData.Data = result
	c.JSON(http.StatusOK, &respData)
	return
}

// 退出
func Logout(c *gin.Context) {
	var respData RespData

	//idStr := c.PostForm("id")
	user, err := global.ParseToken(c.Request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, &respData)
		return
	}

	err = global.Rds.Del(strconv.Itoa(int(user.Id))).Err()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, &respData)
		return
	}

	respData.Data = "成功退出"
	c.JSON(http.StatusOK, &respData)
	return
}

/*
// 获取账户信息
func UserList(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		userListOne(c, id)
	} else {
		userListData(c)
	}
}

// 获取账户信息 - 一条
func userListOne(c *gin.Context, idStr string) {
	var respData RespData
	//验证
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}
	//查询
	var user model.User
	user.Id = id
	userOne, err := user.One()
	if err != nil || userOne.Status == 3 {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}
	userOne.Password = ""

	respData.Data = userOne
	c.JSON(http.StatusOK, &respData)
	return
}

// 获取账户信息 - 列表
func userListData(c *gin.Context) {
	var respData RespData
	// 接收参数
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	// 验证
	if offsetStr != "" {
		_, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			respData.Code = 6000
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	if limitStr != "" {
		_, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			respData.Code = 6000
			c.JSON(http.StatusOK, &respData)
			return
		}
	}

	var user model.User
	var users []model.User
	var err error
	if offsetStr == "" || limitStr == "" {
		users, err = user.All()
	} else {
		users, err = user.Page(offsetStr, limitStr)
	}
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}
	total, err := user.Total()
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	result := struct {
		Total int64        `json:"total"`
		Data  []model.User `json:"data"`
	}{
		Total: total,
		Data:  users,
	}

	respData.Data = result
	c.JSON(http.StatusOK, &respData)
	return
}

// 用户(管理员)修改个人信息
func UserUpdate(c *gin.Context) {
	var respData RespData
	// 接收参数
	idStr := c.PostForm("id")
	name := c.PostForm("name")
	genderStr := c.PostForm("gender")
	phoneStr := c.PostForm("phone")
	identityIdStr := c.PostForm("identityId")
	// 验证
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = filter.MSet(
		filter.El(&name,
			filter.FromString(name, "姓名").
				Required().MinLength(4).MaxLength(16).IsChinese(),
		),
	)
	if err != nil {
		respData.Code = 1013
		c.JSON(http.StatusOK, &respData)
		return
	}
	gender, err := strconv.ParseInt(genderStr, 10, 64)
	if err != nil || (gender != 0 && gender != 1 && gender != 2) {
		respData.Code = 1014
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = global.VerifyMobileFormat(phoneStr)
	if err != nil {
		respData.Code = 1015
		c.JSON(http.StatusOK, &respData)
		return
	}
	phone, err := strconv.ParseInt(phoneStr, 10, 64)
	if err != nil {
		respData.Code = 1015
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = global.IsValidCitizenNo18(identityIdStr)
	if err != nil {
		respData.Code = 1016
		c.JSON(http.StatusOK, &respData)
		return
	}

	var user model.User
	user.Id = id
	userOne, err := user.One()
	if err != nil || userOne.Status == 3 {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 修改
	userOne.Name = name
	userOne.Gender = gender
	userOne.Phone = phone
	userOne.IdentityId = identityIdStr
	err = userOne.Edit()
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	c.JSON(http.StatusOK, &respData)
	return
}

// 用户修改密码
func UserPass(c *gin.Context) {
	var respData RespData
	// 接收参数
	idStr := c.PostForm("id")
	password := c.PostForm("password")
	newpass := c.PostForm("newpass")
	//验证
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = filter.MSet(
		filter.El(&password,
			filter.FromString(password, "旧密码").
				Required().MinLength(6).MaxLength(32).IsLetterOrDigit(),
		),
	)
	if err != nil {
		respData.Code = 1012
		c.JSON(http.StatusOK, &respData)
		return
	}
	err = filter.MSet(
		filter.El(&newpass,
			filter.FromString(newpass, "新密码").
				Required().MinLength(6).MaxLength(32).IsLetterOrDigit(),
		),
	)
	if err != nil {
		respData.Code = 1012
		c.JSON(http.StatusOK, &respData)
		return
	}
	var user model.User
	user.Id = id
	userOne, err := user.One()
	if err != nil || userOne.Status == 3 {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 旧密码验证
	Bool := global.CheckPasswordHash(password, userOne.Password)
	if !Bool {
		respData.Code = 1021
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 新密码加密
	hashpass, err := global.HashPassword(newpass)
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 修改密码
	userOne.Password = hashpass
	err = userOne.Edit()
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	c.JSON(http.StatusOK, &respData)
	return
}

// 管理员修改用户状态
func UserStatus(c *gin.Context) {
	var respData RespData
	// 接收参数
	idStr := c.PostForm("id")
	statusStr := c.PostForm("status")
	// 验证
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}

	status, err := strconv.ParseInt(statusStr, 10, 64)
	if err != nil {
		respData.Code = 1017
		c.JSON(http.StatusOK, &respData)
		return
	}
	if status < 0 || status > 3 {
		respData.Code = 1017
		c.JSON(http.StatusOK, &respData)
		return
	}
	var user model.User
	user.Id = id
	userOne, err := user.One()
	if err != nil || userOne.Status == 3 {
		respData.Code = 1020
		c.JSON(http.StatusOK, &respData)
		return
	}

	userOne.Status = status
	err = userOne.Edit()
	if err != nil {
		log.Println(err)
		respData.Code = 1000
		c.JSON(http.StatusOK, &respData)
		return
	}

	c.JSON(http.StatusOK, &respData)
	return
}
*/