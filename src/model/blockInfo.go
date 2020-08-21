package model

import (
	"fmt"
	"log"

	"src/global"
)

type BlockInfo struct {
	Id     int64  `json:"id"`
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
	Time   string `json:"time"`
}

func (self BlockInfo) GetDataByHeight() (data BlockInfo, err error) {
	sql := "select `id`,`height`,`hash`,`time` from `block_info` where `height`=" + fmt.Sprint(self.Height)
	err = global.Db.QueryRow(sql).Scan(&data.Id, &data.Height, &data.Hash, &data.Time)
	if err != nil {
		log.Println(err)
	}
	return
}

func (self BlockInfo) MaxHeight() (maxHeight int64, err error) {
	sql := "select max(`height`) as `maxHeight` from `block_info`"
	err = global.Db.QueryRow(sql).Scan(&maxHeight)
	if err != nil {
		log.Println(err)
	}
	return
}

func (self BlockInfo) Add(str string) (err error) {
	var sql string = "insert into `block_info` (`height`,`hash`,`time`) value " + str
	_, err = global.Db.Exec(sql)
	if err != nil {
		log.Println(err)
	}
	return
}

//func (self BlockInfo) Add() (err error) {
//	var sql string = "insert into `block_info` (`height`,`hash`,`time`) value (?,?,?)"
//	_, err = global.Db.Exec(sql, self.Height, self.Hash, self.Time)
//	if err != nil {
//		log.Println(err)
//	}
//	return
//}
