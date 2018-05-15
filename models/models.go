package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//注册数据库
func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:@/tzh?charset=utf8")
	orm.RegisterModel()
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default
}

//执行sql语句
func Query(a string) (maps []orm.Params) {
	orm.Debug = false
	o := orm.NewOrm()
	_, err := o.Raw(a).Values(&maps)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//执行 update delete操作
func Exec(a string) error {
	orm.Debug = false
	o := orm.NewOrm()
	_, err := o.Raw(a).Exec()
	if err != nil {
		fmt.Println("exec bug:", err)
	}
	return err
}
