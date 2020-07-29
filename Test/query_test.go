package Test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yoyofxteam/yoyodata/cache"
	"github.com/yoyofxteam/yoyodata/queryable"
	"testing"
)

type User struct {
	Id         int
	UserName   string
	Age        int
	Department string
}

func Test(t *testing.T) {
	cache.NewTypeInfoCache()
	query := queryable.Queryable{
		DB: queryable.DbInfo{
			DriverName: "mysql",
			DataSource: "root:A.jiheMA?1@tcp(49.232.153.51)/go_study",
		},
		Model: User{},
	}
	var userArray []User
	query.Query("select age,username from t_user", &userArray)
	fmt.Print(userArray)
}

func Test2(t *testing.T) {
	/*cache.NewTypeInfoCache()
	query := queryable.Queryable{
		DB: queryable.DbInfo{
			DriverName: "mysql",
			DataSource: "root:A.jiheMA?1@tcp(49.232.153.51)/go_study",
		},
		Model: User{},
	}
	var userArray []User
	query.QueryByParams("select age,username from t_user WHERE age=?", &userArray,23)
	fmt.Print(userArray)*/

	var strss= []string{
		"qwr",
		"234",
		"yui",
		"cvbc",
	}
	test(strss...)
}

func  test(str ...string)  {
	fmt.Println(str)
}


