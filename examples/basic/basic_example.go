package main

import (
	"time"

	"fmt"

	"github.com/JessonChan/yorm"
)

type ProgramLanguage struct {
	Id        int64
	Position  int
	Name      string
	RankMonth time.Time
	Created   time.Time
}

func main() {
	yorm.SetLoggerLevel(yorm.DebugLevel)
	yorm.Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")

	//插入一条数据
	php := ProgramLanguage{Name: "PHP", Position: 7, RankMonth: time.Now(), Created: time.Now()}
	yorm.Insert(&php)

	var ps []ProgramLanguage

	//读取所有的数据
	yorm.R(&ps)
	fmt.Println(ps)

	//读取所有小于10的数据
	yorm.R(&ps, "where id<10")
	fmt.Println(ps)

	//也可以
	yorm.Select(&ps, "where <10")
	fmt.Println(ps)

	var p ProgramLanguage = ProgramLanguage{Id: 1}

	//读取id为1的某条数据
	yorm.R(&p)
	fmt.Println(p)

	//读取id为1的某条数据
	yorm.SelectByPK(&p)
	fmt.Println(p)

	//读取id为2的某条数据,
	yorm.SelectByPK(&p, "where id=?", 2)
	fmt.Println(p)

	//更新一条数据
	yorm.Update("update program_language set position=?", 8)

	//删除一条
	yorm.Delete("delete from program_language where id=? ", p.Id)
}
