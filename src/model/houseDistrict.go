package model

import (
	"fmt"

	"src/global"
)

type Housedistrict struct {
	Id         int64  `json:"id"`
	Areaid     int64  `json:"areaid"`
	Coordinate string `json:"coordinate"`
	Attribute  int64  `json:"attribute"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
}

// 获取所有数据或指定attribute属性的所有数据
func (self Housedistrict) All(attr string) (data []Housedistrict, err error) {
	sql := "SELECT `id`,`areaid`,`coordinate`,`attribute`,`name`,`icon` FROM `house_district`;"
	if attr != "" {
		sql = fmt.Sprintf("SELECT `id`,`areaid`,`coordinate`,`attribute`,`name`,`icon` FROM `house_district` WHERE `attribute`=%v;", fmt.Sprint(attr))
	}
	rows, err := global.Db.Query(sql)
	defer rows.Close()
	if err != nil {
		return
	}
	var tmp Housedistrict
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Areaid, &tmp.Coordinate, &tmp.Attribute, &tmp.Name, &tmp.Icon)
		if err != nil {
			return
		}
		data = append(data, tmp)
	}
	return
}
