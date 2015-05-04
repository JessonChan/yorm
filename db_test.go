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

func TestYorm(t *testing.T) {
	err := Register("root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	p := ProgramLanguage{Name: "PHP", Position: 7, RankMonth: time.Now(), Created: time.Now()}
	_, err = Insert(&p)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(p)
	Update("update program_language set position=? where id=? ", 12, p.Id)

	var p1 ProgramLanguage
	Select(&p1, "select * from program_language where id=?", p.Id)

	var p2 ProgramLanguage
	Select(&p2, "where id=?", p.Id)
	if p2.Id != p.Id {
		t.Log(p2)
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
}
