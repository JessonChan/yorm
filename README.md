# README #

yOrm is a simple,lightweight orm  , for mysql only now.

### Why this project calls yOrm ###

yOrm is just a name.
more about the detail, cc [https://github.com/lewgun]

### What is this yOrm for? ###

* A simple mysql orm to crud

## Tags ##
 
Now support these types of tag.

### column ###
this tag alias struct name to a real column name. "Id int \`yorm:column(autoId)\`" means this field Id will name autoId in mysql column

### pk ###
this tag allow you to set a primary key where select/delete/update as the where clause  "Id int \`yorm:column(autoId);pk\`"


# benchmark #

select by id with five fields

> beegoOrm 13376 microsecond   
>     xorm 16718 microsecond   
>     yorm 6759 microsecond   

code is here:



<pre>
<code>


package main

import (
	"beego/orm"
	"fmt"
	"time"

	"github.com/JessonChan/fastfunc"
	"github.com/JessonChan/yorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const db = `root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8`

type ProgramLanguage struct {
	Id        int64
	Name      string
	RankMonth time.Time
	Position  int
	Created   time.Time
}

var engine *xorm.Engine
var o orm.Ormer

func init() {
	orm.RegisterDataBase("default", "mysql", db)
	orm.RegisterModel(new(ProgramLanguage))
	yorm.Register(db)
	engine, _ = xorm.NewEngine("mysql", db)
	o = orm.NewOrm()
}

func main() {
	fastfunc.SetRunTimes(1e5)
	fmt.Println("beegoOrm", fastfunc.Run(beegoOrm)/1e6, "microsecond")
	fmt.Println("    xorm", fastfunc.Run(xomrTest)/1e6, "microsecond")
	fmt.Println("    yorm", fastfunc.Run(yormTest)/1e6, "microsecond")
}

func beegoOrm() {
	p := ProgramLanguage{Id: 1}
	o.Read(&p)
	if p.Name == "" {
		panic(p)
	}
}
func yormTest() {
	p := ProgramLanguage{Id: 1}
	yorm.SelectByPK(&p)
	if p.Name == "" {
		panic(p)
	}
}
func xomrTest() {
	p := ProgramLanguage{Id: 1}
	engine.Get(&p)
	if p.Name == "" {
		panic(p)
	}
}

</code>
</pre>
