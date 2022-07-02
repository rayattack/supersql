package supersql_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/rayattack/supersql"
)

const DSN = "postgres://tester:tester@localhost:5432/dvdrental?sslmode=disable"
const (
	age   = "age"
	name  = "name"
	email = "email"
)

var Xql *supersql.SqlQuery

func init() {
	sql, err := supersql.Query(context.Background(), DSN)
	if err != nil {
		log.Fatalf("connection failed: %s", err)
	}
	Xql = sql
}

func TestQueryConnectionClose(t *testing.T) {
	conn, err := supersql.Query(context.Background(), DSN)
	if err != nil {
		t.FailNow()
	} else {
		err := conn.CLOSE()
		if err != nil {
			t.Logf("this is the error: %s", err)
		}
	}
}

func TestCreateCommand(t *testing.T) {
}

func TestExecCommand(t *testing.T) {
	q := Xql.SELECT("title").FROM("film").WHERE("film_id = ?", 133)
	r, err := q.GO()
	if err != nil {
		t.Fail()
		t.Logf("error occured: %s", err)
	}
	_, ok := r.(supersql.SqlResult)
	if !ok {
		t.Fail()
	}
}

func TestJoinCommand(t *testing.T) {
	//use sql -- comment operator to mark where spaces will be added after \t and \n removed
	sql := `
		SELECT c.email, c.first_name || ', ' || c.last_name AS name, r.rental_date, f.title--
			FROM rental r--
		JOIN inventory i ON--
			i.inventory_id = r.inventory_id--
		JOIN film f ON--
			f.film_id = i.film_id--
		JOIN customer c ON--
			c.customer_id = r.customer_id
	`

	//reformat sql string above removing and tabs and spaces so
	//it can be matched with expantiated query afterwards
	sql = strings.Replace(sql, "\t", "", -1)
	sql = strings.Replace(sql, "\n", "", -1)
	sql = strings.Replace(sql, "--", " ", -1)

	q := generateJoinCommand()
	if q.PP() != sql {
		t.Log(sql)
		t.Log(q.PP())
		t.Fail()
	}
}

func TestSelectCommand(t *testing.T) {
	q := Xql.SELECT("first")
	if q.PP() != "SELECT first" {
		t.Fail()
	}

	r := q.FROM("film")
	if r.PP() != "SELECT first FROM film" {
		t.Fail()
	}

	s := Xql.SELECT("title").FROM("film").WHERE("film_id = ?", 133)
	if s.PP() != "SELECT title FROM film WHERE film_id = 133" {
		t.Log(s.PP())
		t.Fail()
	}
}

func TestLimitCommand(t *testing.T) {
	q := Xql.SELECT().FROM(actor).LIMIT(5)
	if q.PP() != "SELECT * FROM actor LIMIT 5" {
		t.Fail()
	}
}

func TestOffsetCommand(t *testing.T) {
	q := Xql.SELECT().FROM(rental).LIMIT(5).OFFSET(4)
	if q.PP() != "SELECT * FROM rental LIMIT 5 OFFSET 4" {
		t.Fail()
	}
}

func TestOrderByAscDesc(t *testing.T) {
	b := Xql.SELECT().FROM(actor)
	ascending := "SELECT * FROM actor ORDER BY actor_id ASC"
	descending := "SELECT * FROM actor ORDER BY actor_id DESC"

	q := b.ORDER_BY("actor_id")
	if q.PP() != "SELECT * FROM actor ORDER BY actor_id" {
		t.Fail()
	}

	q = b.ORDER_BY("actor_id ASC")
	if q.PP() != ascending {
		t.Fail()
		t.Log(q.PP())
	}
	q = b.ORDER_BY("actor_id").ASC()
	if q.PP() != ascending {
		t.Fail()
	}
	q = b.ASC("actor_id")
	if q.PP() != ascending {
		t.Fail()
	}

	q = b.DESC("actor_id")
	if q.PP() != descending {
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	q := Xql.INSERT(name, age, email).INTO(actor)
	if q.PP() != "INSERT INTO actor (name, age, email)" {
		t.Fail()
	}

	old := q.PP()

	vf := func(vals ...interface{}) []interface{} {
		return vals
	}
	q = q.VALUES(vf("Ryan Bryan", 4, "mc@ortserga.com"))
	if q.PP() != fmt.Sprintf("%s VALUES (?, ?, ?)", old) {
		t.Fail()
	}
}

func TestInsertInto(t *testing.T) {
	s := Xql.INSERT_INTO("customers (?, ?, ?)", []string{name, age, email})
	q := Xql.INSERT_INTO(actor, []string{name, age, email}).VALUES([]interface{}{"ab", 4, "cdef"})
	l := Xql.INSERT_INTO(actor, []string{"name", "nom"}) //consider changing this to slice instead???
	x := Xql.INSERT(name, age, email).INTO(actor).VALUES([]interface{}{"ab", 5, "cdef"})

	if x.PP() != q.PP() {
		t.Fail()
	}
	if s.PP() != "INSERT INTO customers (name, age, email)" {
		t.Fail()
	}
	if q.PP() != "INSERT INTO actor (name, age, email) VALUES (?, ?, ?)" {
		t.Fail()
	}
	if l.PP() != "INSERT INTO actor (name, nom)" {
		t.Fail()
	}
}
