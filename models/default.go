package models

import (
	"time"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Prefix = beego.AppConfig.String("database::db_prefix")
	DbType = beego.AppConfig.String("database::db_type")
	DbUrl  = beego.AppConfig.String("database::db_url")
)

func init() {
	orm.Debug, _ = beego.AppConfig.Bool("orm.debug")
	orm.DefaultTimeLoc = time.UTC

	orm.RegisterModelWithPrefix(Prefix, new(BuildList))
	orm.RegisterModelWithPrefix(Prefix, new(BuildListPkg))
	orm.RegisterModelWithPrefix(Prefix, new(Karma))

	orm.RegisterModelWithPrefix(Prefix, new(Advisory))

	orm.RegisterModelWithPrefix(Prefix, new(User))
	orm.RegisterModelWithPrefix(Prefix, new(UserPermission))

	err := orm.RegisterDataBase("default", DbType, DbUrl, 30)
	if err != nil {
		panic(err)
	}

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
}
