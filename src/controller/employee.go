package controller

import (
	"log"
	"net/http"
	"src/model"
	"strconv"

	"github.com/dxvgef/filter"

	"github.com/gin-gonic/gin"
)

func EmployeeList(c *gin.Context) {
	employeeid := c.Query("id")
	if employeeid == "" {
		employeePage(c)
	} else {
		employeeOne(c, employeeid)
	}
}

// 分页
func employeePage(c *gin.Context) {
	var respData RespData
	// 接收数据
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	ridStr := c.Query("rid")
	// 验证
	var err error
	//var offset int64
	if offsetStr != "" {
		_, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			respData.Error = "偏移量必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	//var limit int64
	if limitStr != "" {
		_, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			respData.Error = "每页显示的条数必须是数字"
			c.JSON(http.StatusOK, &respData)
			return
		}
	}
	var rid int64
	if ridStr != "" {
		rid, err = strconv.ParseInt(ridStr, 10, 64)
		if err != nil || rid <= 0 {
			respData.Error = "餐厅ID必须是数字，并且大于0"
			c.JSON(http.StatusOK, &respData)
			return
		}
		var rest model.Restaurnt
		rest.Id = rid
		_, err := rest.One()
		if err != nil {
			respData.Error = "餐厅ID不存在"
			c.JSON(http.StatusBadRequest, &respData)
			return
		}
	}

	// 获取数据
	var employee model.Employee
	var data []model.Employee
	if offsetStr == "" || limitStr == "" {
		data, err = employee.GetDataByRid(rid)
	} else {
		data, err = employee.Page(rid, offsetStr, limitStr)
	}
	if err != nil {
		log.Println(err)
		respData.Error = "获取设备信息列表失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	// 获取数据总数
	num, err := employee.Total(rid)
	if err != nil {
		log.Println(err)
		respData.Error = "获取设备信息总条数失败"
		c.JSON(http.StatusOK, &respData)
		return
	}

	res := struct {
		Num  int64            `json:"num"`
		Data []model.Employee `json:"data"`
	}{
		Num:  num,
		Data: data,
	}

	respData.Data = res
	c.JSON(http.StatusOK, &respData)
	return

}

//获取一条设备信息数据
func employeeOne(c *gin.Context, employeeidStr string) {
	var respData RespData
	// 接收数据
	var employee model.Employee
	// 验证数字类型
	employeeid, err := strconv.ParseInt(employeeidStr, 10, 64)
	if err != nil {
		respData.Error = "员工ID必须是数字"
		c.JSON(http.StatusOK, &respData)
		return
	}
	// 通过ID查询信息
	employee.Id = employeeid
	data, err := employee.One()
	if err != nil {
		respData.Error = "员工不存在"
		c.JSON(http.StatusOK, &respData)
		return
	}

	respData.Data = data
	c.JSON(http.StatusOK, &respData)
	return
}

func EmployeeAdd(c *gin.Context) {
	var respData RespData

	// 接收数据
	name := c.Request.FormValue("name")
	sexStr := c.Request.FormValue("sex")
	ageStr := c.Request.FormValue("age")
	position := c.Request.FormValue("position")
	onDutyStr := c.Request.FormValue("onDuty")
	healthTestStr := c.Request.FormValue("healthTest")
	timeTestStr := c.Request.FormValue("timeTest")
	temperatureStr := c.Request.FormValue("temperature")
	ridStr := c.Request.FormValue("rid")

	err := filter.MSet(
		filter.El(&name,
			filter.FromString(name, "员工名称").
				Required().MaxLength(32).IsLetterOrDigit("员工名称不能含有特殊字符"),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	sex, err := strconv.ParseInt(sexStr, 10, 64)
	if err != nil || (sex != 0 && sex != 1) {
		respData.Error = "填写正确的性别参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	age, err := strconv.ParseInt(ageStr, 10, 64)
	if err != nil || age < 18 || age > 80 {
		respData.Error = "填写正确的年龄参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	err = filter.MSet(
		filter.El(&position,
			filter.FromString(position, "职位名称").
				Required().MaxLength(64).IsLetterOrDigit("职位名称不能含有特殊字符"),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	onDuty, err := strconv.ParseInt(onDutyStr, 10, 64)
	if err != nil || (onDuty != 0 && onDuty != 1) {
		respData.Error = "填写正确的在岗参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	healthTest, err := strconv.ParseInt(healthTestStr, 10, 64)
	if err != nil || (healthTest != 0 && healthTest != 1) {
		respData.Error = "填写正确的健康检测结果参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	timeTest, err := strconv.ParseInt(timeTestStr, 10, 64)
	if err != nil || timeTest < 0 {
		respData.Error = "填写正确的检测时间参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil || temperature <= 30 {
		respData.Error = "填写正确的体温参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil || rid <= 0 {
		respData.Error = "填写正确的餐厅ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var rest model.Restaurnt
	rest.Id = rid
	_, err = rest.One()
	if err != nil {
		respData.Error = "餐厅ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	var employee model.Employee
	employee.Rid = rid
	employee.Name = name
	employee.Sex = sex
	employee.Age = age
	employee.Position = position
	employee.OnDuty = onDuty
	employee.HealthTest = healthTest
	employee.TimeTest = timeTest
	employee.Temperature = temperature
	employee.Status = 1
	err = employee.Add()
	if err != nil {
		respData.Error = "添加失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "添加成功"
	c.JSON(http.StatusOK, &respData)
	return
}

func EmployeeEdit(c *gin.Context) {
	var respData RespData

	// 接收数据
	name := c.Request.FormValue("name")
	sexStr := c.Request.FormValue("sex")
	ageStr := c.Request.FormValue("age")
	position := c.Request.FormValue("position")
	onDutyStr := c.Request.FormValue("onDuty")
	healthTestStr := c.Request.FormValue("healthTest")
	timeTestStr := c.Request.FormValue("timeTest")
	temperatureStr := c.Request.FormValue("temperature")
	ridStr := c.Request.FormValue("rid")
	employeeidStr := c.Request.FormValue("id")

	err := filter.MSet(
		filter.El(&name,
			filter.FromString(name, "员工名称").
				Required().MaxLength(32).IsLetterOrDigit("员工名称不能含有特殊字符"),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	sex, err := strconv.ParseInt(sexStr, 10, 64)
	if err != nil || (sex != 0 && sex != 1) {
		respData.Error = "填写正确的性别参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	age, err := strconv.ParseInt(ageStr, 10, 64)
	if err != nil || age < 18 || age > 80 {
		respData.Error = "填写正确的年龄参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	err = filter.MSet(
		filter.El(&position,
			filter.FromString(position, "职位名称").
				Required().MaxLength(64).IsLetterOrDigit("职位名称不能含有特殊字符"),
		),
	)
	if err != nil {
		respData.Error = err.Error()
		c.JSON(http.StatusOK, &respData)
		return
	}

	onDuty, err := strconv.ParseInt(onDutyStr, 10, 64)
	if err != nil || (onDuty != 0 && onDuty != 1) {
		respData.Error = "填写正确的在岗参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	healthTest, err := strconv.ParseInt(healthTestStr, 10, 64)
	if err != nil || (healthTest != 0 && healthTest != 1) {
		respData.Error = "填写正确的健康检测结果参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	timeTest, err := strconv.ParseInt(timeTestStr, 10, 64)
	if err != nil || timeTest < 0 {
		respData.Error = "填写正确的检测时间参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil || temperature <= 30 {
		respData.Error = "填写正确的体温参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}

	rid, err := strconv.ParseInt(ridStr, 10, 64)
	if err != nil || rid <= 0 {
		respData.Error = "填写正确的餐厅ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var rest model.Restaurnt
	rest.Id = rid
	_, err = rest.One()
	if err != nil {
		respData.Error = "餐厅ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	employeeid, err := strconv.ParseInt(employeeidStr, 10, 64)
	if err != nil || employeeid <= 0 {
		respData.Error = "填写正确的员工ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var employee model.Employee
	employee.Id = employeeid
	one, err := employee.One()
	if err != nil {
		respData.Error = "员工ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	one.Rid = rid
	one.Name = name
	one.Sex = sex
	one.Age = age
	one.Position = position
	one.OnDuty = onDuty
	one.HealthTest = healthTest
	one.TimeTest = timeTest
	one.Temperature = temperature
	err = one.Edit()
	if err != nil {
		respData.Error = "修改失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "修改成功"
	c.JSON(http.StatusOK, &respData)
	return
}

func EmployeeDel(c *gin.Context) {
	var respData RespData

	employeeidStr := c.Request.FormValue("id")

	employeeid, err := strconv.ParseInt(employeeidStr, 10, 64)
	if err != nil || employeeid <= 0 {
		respData.Error = "填写正确的员工ID参数！"
		c.JSON(http.StatusOK, &respData)
		return
	}
	var employee model.Employee
	employee.Id = employeeid
	one, err := employee.One()
	if err != nil {
		respData.Error = "员工ID不存在"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	one.Status = 0
	err = one.Edit()
	if err != nil {
		respData.Error = "删除失败"
		c.JSON(http.StatusBadRequest, &respData)
		return
	}

	respData.Data = "删除成功"
	c.JSON(http.StatusOK, &respData)
	return
}
