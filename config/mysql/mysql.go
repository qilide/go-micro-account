package mysql

import (
	"github.com/jinzhu/gorm"
)

// MysqlInit 初始化数据库
func MysqlInit(user string,pwd string,database string)(*gorm.DB,error){
	//连接数据库
	db, err := gorm.Open("mysql", user+":"+pwd+"@/"+database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return db,err
	}
	//禁止复表
	db.SingularTable(true)
	//初始化表，只使用第一次
	//err = repository.NewUserRepository(db).InitTable()
	//if err != nil {
	//	logger.Error(err)
	//}
	return db,nil
}
