package yorm

import (
	"testing"
	"time"
)

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ProgramLanguage struct {
	Id        int64
	Position  int
	Name      string
	RankMonth time.Time
	Created   time.Time
}
type GolangWord struct {
	Aid  int `yorm:"pk"`
	Word string
	Rate float32
}

type A struct {
	Aid  int `yorm:"pk"`
	G2   GolangWord
	P1   ProgramLanguage
	Word string
	Rate float64
	GW   GolangWord
}

func TestCount(t *testing.T) {
	SetLoggerLevel(Debug)
	Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	//	t.Log(Count(&GolangWord{}))
	var c int64
	c = 64
	Select(&c, "select count(0) from program_language where id>100000")
	if c != 0 {
		t.FailNow()
	}
	t.Log(Count(&GolangWord{}, "where rate>?", 0.5))
}

func TestTranSelect(t *testing.T) {
	SetLoggerLevel(Debug)
	Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	p := ProgramLanguage{Name: "PHP", Position: 7, RankMonth: time.Now(), Created: time.Now()}
	tran, err := Begin()
	fmt.Println(err)
	tran.Insert(&p)
	id, err := Insert(&p)
	if p.Id > 0 {
		fmt.Println(RollBack(tran))
	}
	fmt.Println(Commit(tran))
	t.Log(id, p.Id)
}

func TestYorm(t *testing.T) {
	SetLoggerLevel(Debug)
	err := Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	RegisterWithName("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8", "read")
	RegisterWithName("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8", "write")
	p := ProgramLanguage{Name: "PHP", Position: 7, RankMonth: time.Now(), Created: time.Now()}
	_, err = Insert(&p)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(p)
	_, err = Using("write").Update("update program_language set position=? where id=? ", 12, p.Id)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var p1 ProgramLanguage
	Using("read").Select(&p1, "select * from program_language where id=?", p.Id)

	var p2 ProgramLanguage
	Select(&p2, "where id=?", p.Id)
	if p2.Id != p.Id {
		t.Log(p2)
		t.FailNow()
	}
	var p3 = ProgramLanguage{Id: p.Id}
	SelectByPK(&p3)
	t.Log(p3)

	if p3.Name != p2.Name {
		t.Log(p3)
		t.FailNow()
	}

	if p1.Name != p.Name {
		t.Log(p1)
		t.FailNow()
	}
	if p1.Position != 12 {
		t.Log(p1)
		t.FailNow()
	}
	t.Log(p1)
	Delete("delete from program_language where id=? ", p.Id)
	err = Select(&p1, "select * from program_language where id=?", p.Id)
	t.Log(err)
	if err == nil {
		t.FailNow()
	}
	err = Using("nil").Select(&p1, "select * from program_language where id=?", p.Id)
	if err != ErrNilSqlExecutor {
		t.FailNow()
	}
	var count int
	Select(&count, "select count(0) from program_language")
	t.Log(count)
	if count == 0 {
		t.FailNow()
	}
}

func TestSelectByPk(t *testing.T) {
	Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	g := GolangWord{Aid: 1}
	err := SelectByPK(&g)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(g)

	a := A{Aid: 2}
	err = SelectByPK(&a, "golang_word")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(a)
}

func TestListSelect(t *testing.T) {
	SetLoggerLevel(Debug)
	Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	var ps []ProgramLanguage
	err := Select(&ps, "select * from program_language")
	if len(ps) == 0 {
		t.Log(err)
		t.FailNow()
	}
	for _, v := range ps {
		t.Log(v)
	}
}

func TestSelectIdList(t *testing.T) {
	Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	var ids []int
	R(&ids, "select id from program_language")
	t.Log(ids)
}

func TestSMethod(t *testing.T) {
	SetLoggerLevel(Debug)
	Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	var ps []ProgramLanguage
	err := R(&ps)
	if len(ps) == 0 {
		t.Log(err)
		t.FailNow()
	}
	for _, v := range ps {
		t.Log(v)
	}
	var p ProgramLanguage = ProgramLanguage{Id: 1}
	err = R(&p)
	t.Log(p)
	if p.Name == "" {
		t.Log(err)
		t.FailNow()
	}

	err = R(&ps, "select * from program_language")
	if len(ps) == 0 {
		t.Log(err)
		t.FailNow()
	}
	for _, v := range ps {
		t.Log(v)
	}
	err = R(&p, "select * from program_language where id=?", 2)
	t.Log(p)
	if p.Name == "" {
		t.Log(err)
		t.FailNow()
	}
}
