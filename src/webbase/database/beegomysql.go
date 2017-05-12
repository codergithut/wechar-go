package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 导入数据库驱动
	"time"
)


type Userinfo struct {
	Uid     int `PK` //如果表的主键不是id，那么需要加上pk注释，显式的说这个字段是主键
	Username    string
	Departname  string
	Created     time.Time
}

type User struct {
	Id          int `PK` //如果表的主键不是id，那么需要加上pk注释，显式的说这个字段是主键
	Name        string
	Profile     *Profile   `orm:"rel(one)"` // OneToOne relation
	Post        []*Post `orm:"reverse(many)"` // 设置一对多的反向关系
}

type Profile struct {
	Id          int
	Age         int16
	User        *User   `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type Post struct {
	Id    int
	Title string
	User  *User  `orm:"rel(fk)"`    //设置一对多关系
	Tags  []*Tag `orm:"rel(m2m)"`
}

type Tag struct {
	Id    int
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}

func init() {
	// 设置默认数据库
	orm.RegisterDataBase("default", "mysql", "root:root@/mysql?charset=utf8", 30)

	// 注册定义的 model
	orm.RegisterModel(new(User))
	//RegisterModel 也可以同时注册多个 model
	//orm.RegisterModel(new(User), new(Profile), new(Post))

	// 创建 table
	orm.RunSyncdb("mysql", false, true)
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Userinfo),new(User), new(Profile), new(Tag), new(Post))

	//根据数据库的别名，设置数据库的最大空闲连接
	orm.SetMaxIdleConns("mysql", 30)

	//根据数据库的别名，设置数据库的最大数据库连接 (go >= 1.2)
	orm.SetMaxOpenConns("mysql", 30)

	//目前beego orm支持打印调试，你可以通过如下的代码实现调试
	orm.Debug = true
}

func main() {
	o := orm.NewOrm()
	var user Userinfo
	user.Username="test"
	user.Departname="test"
	user.Uid=123;

	id, err := o.Insert(&user)
	if err == nil {
		fmt.Println(id)
	}
}