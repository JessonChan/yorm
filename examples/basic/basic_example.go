package main

import (
	"time"

	"fmt"

	"github.com/JessonChan/yorm"
)

/*
CREATE TABLE `program_language` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) DEFAULT NULL,
  `rank_month` date DEFAULT NULL,
  `position` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;
*/

type ProgramLanguage struct {
	Id        int64
	Position  int
	Name      string
	RankMonth time.Time
	Created   time.Time
}

func main() {
	yorm.SetLoggerLevel(yorm.DebugLevel)
	//设置自己的数据地址
	yorm.Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")

	//插入一条数据
	php := ProgramLanguage{Name: "PHP", Position: 7, RankMonth: time.Now(), Created: time.Now()}
	yorm.Insert(&php)

	//更新一条数据
	yorm.Update("update program_language set position=? where id=?", 8, php.Id)

	//删除一条
	yorm.Delete("delete from program_language where id=? ", php.Id)

	var ps []ProgramLanguage

	//读取所有的数据
	yorm.R(&ps)
	fmt.Println(ps)

	//读取所有小于10的数据
	yorm.R(&ps, "where id<10")
	fmt.Println(ps)

	//读取所有小于10的数据
	yorm.R(&ps, "where id<?", 10)
	fmt.Println(ps)

	//也可以
	yorm.Select(&ps, "where id<?", 10)
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

}
