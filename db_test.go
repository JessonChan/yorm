package yorm

import (
	"testing"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

type ProgramLanguage struct {
	Id        int64
	Name      string
	RankMonth time.Time
	Position  int
	Created   time.Time
}
type GolangWord struct {
	Aid  int `yorm:"pk"`
	Word string
	Rate float32
}

type A struct {
	Aid  int `yorm:"pk"`
	Word string
	Rate float64
}

func TestYorm(t *testing.T) {
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
	Using("write").Update("update program_language set position=? where id=? ", 12, p.Id)

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


