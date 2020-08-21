package model

import (
	"fmt"
	"log"

	"src/global"
)

type StatusType struct {
	Id         int64  `json:"id"`
	Statusid   string `json:"statusid"`
	StatusName string `json:"statusname"`
	Parentid   string `json:"parentid"`
	Brand      string `json:"brand"`
	Des        string `json:"des"`
	Status     int64  `json:"status"`
}

//获取所有device type 数据
func (self StatusType) All() (data []StatusType, err error) {
	var sql string = "select `id`,`statusid`,`statusname`,`parentid`,`brand`,`des`,`status` from `status_type` where `status`=1 "
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp StatusType
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Statusid, &tmp.StatusName, &tmp.Parentid, &tmp.Brand, &tmp.Des, &tmp.Status)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

func (self StatusType) Pages(offset int64, limit int64, keys string) (data []StatusType, err error) {
	var sql string = fmt.Sprintf("select `id`,`statusid`,`statusname`,`parentid`,`brand`,`des`,`status` from `status_type` where `status`=1 order by `id` desc limit %v, %v", offset, limit)
	if keys != "" {
		sql = fmt.Sprintf("select `id`,`statusid`,`statusname`,`parentid`,`brand`,`des`,`status` from `status_type` where `status`=1 and (`statusname` like '%v' or `des` like '%v') order by `id` desc limit %v, %v", "%"+keys+"%", "%"+keys+"%", offset, limit)
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp StatusType
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Statusid, &tmp.StatusName, &tmp.Parentid, &tmp.Brand, &tmp.Des, &tmp.Status)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, tmp)
	}
	return
}

//获取一条device type数据
func (self StatusType) One() (data StatusType, err error) {
	var sql string = fmt.Sprintf("select `id`,`statusid`,`statusname`,`parentid`,`brand`,`des`,`status` from `status_type` where `status`=1 and `id`='%v'", self.Id)
	err = global.Db.QueryRow(sql).Scan(&data.Id, &data.Statusid, &data.StatusName, &data.Parentid, &data.Des, &data.Des, &data.Status)
	return
}

// 获取报警ids
func (self StatusType) GetAlarmIds() (AlarmIds []int64, err error) {
	sql := "select `id` from `status_type` where `status`=1 and `parentid`='0x1'"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp int64
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		AlarmIds = append(AlarmIds, tmp)
	}
	return
}

// 获取异常ids
func (self StatusType) GetExceptionIds() (exIds []int64, err error) {
	sql := "select `id` from `status_type` where `status`=1 and `parentid`='0x2'"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp int64
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		exIds = append(exIds, tmp)
	}
	return
}

// 获取操作ids
func (self StatusType) GetOperationIds() (operationIds []int64, err error) {
	sql := "select `id` from `status_type` where `status`=1 and `parentid`='0x3'"
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	var tmp int64
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		operationIds = append(operationIds, tmp)
	}
	return
}

// 根据parentid获取数据
func (self StatusType) GetDataByParentid() (operationIds []int64, err error) {
	sql := fmt.Sprintf("select `id` from `status_type` where `status`=1 and `parentid`='%v'", self.Parentid)
	rows, err := global.Db.Query(sql)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	var tmp int64
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			log.Println(err)
			return
		}
		operationIds = append(operationIds, tmp)
	}
	return
}
