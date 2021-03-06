package model

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"

	. "github.com/easykoo/go-display/common"

	"os"
	"time"
)

var orm *xorm.Engine

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func SetEngine() *xorm.Engine {
	Log.Info("db initializing...")
	var err error
	orm, err = xorm.NewEngine("sqlite3", "./.db")
	PanicIf(err)
	orm.TZLocation = time.Local
	orm.ShowSQL = Cfg.MustBool("db", "show_sql", false)
	orm.Logger = xorm.NewSimpleLogger(Log.GetWriter())
	InitDB()
	Log.Info("db initialized...")
	return orm
}

func InitDB() {
	if orm == nil {
		Log.Err("db not exists")
	}

	if exist, _ := orm.IsTableExist(&Pad{}); !exist {
		orm.DropTables(&Pad{}, &Picture{})
	}

	if exist, _ := orm.IsTableExist(&Module{}); !exist {
		orm.CreateTables(&Module{})
		orm.InsertOne(&Module{Id: 1, Description: "Admin"})
		orm.InsertOne(&Module{Id: 2, Description: "User"})
		orm.InsertOne(&Module{Id: 3, Description: "Pad"})
		orm.InsertOne(&Module{Id: 4, Description: "Picture"})
	}
	if exist, _ := orm.IsTableExist(&Privilege{}); !exist {
		orm.CreateTables(&Privilege{})
		orm.InsertOne(&Privilege{ModuleId: 1, RoleId: 1, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 1, RoleId: 2, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 2, RoleId: 1, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 2, RoleId: 2, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 2, RoleId: 3, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 3, RoleId: 1, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 3, RoleId: 2, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 3, RoleId: 3, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 4, RoleId: 1, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 4, RoleId: 2, DeptId: 1})
		orm.InsertOne(&Privilege{ModuleId: 4, RoleId: 3, DeptId: 1})
	}
	if exist, _ := orm.IsTableExist(&Dept{}); !exist {
		orm.CreateTables(&Dept{})
		orm.InsertOne(&Dept{Id: 1, Description: "Default"})
	}
	if exist, _ := orm.IsTableExist(&Role{}); !exist {
		orm.CreateTables(&Role{})
		orm.InsertOne(&Role{Id: 1, Description: "Admin"})
		orm.InsertOne(&Role{Id: 2, Description: "Manager"})
		orm.InsertOne(&Role{Id: 3, Description: "Employee"})
		orm.InsertOne(&Role{Id: 4, Description: "User"})
	}
	if exist, _ := orm.IsTableExist(&User{}); !exist {
		orm.CreateTables(&User{})
		orm.InsertOne(&User{Username: "admin", Password: "b0baee9d279d34fa1dfd71aadb908c3f", FullName: "Admin", Role: Role{Id: 1}, Dept: Dept{Id: 1}})
	}
	if exist, _ := orm.IsTableExist(&Settings{}); !exist {
		orm.CreateTables(&Settings{})
		orm.InsertOne(&Settings{Id: 1, AppName: "Display", Owner: User{Id: 1}})
	}
	if exist, _ := orm.IsTableExist(&Pad{}); !exist {
		orm.CreateTables(&Pad{})
	}
	if exist, _ := orm.IsTableExist(&Picture{}); !exist {
		orm.CreateTables(&Picture{})
	}
	if exist, _ := orm.IsTableExist(&PadPicture{}); !exist {
		orm.CreateTables(&PadPicture{})
	}
}
