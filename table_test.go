package supersql_test

import (
	"testing"

	"github.com/rayattack/supersql"
)

var (
	film = supersql.Table("film")
	rental supersql.Relation = supersql.Table("rental")
)

func TestTableHelper(t *testing.T){
	query := Xql.SELECT().FROM(film)
	if query.PP() != "SELECT * FROM film" {
		t.Fail()
	}

	rental.AS("r")
	query = Xql.SELECT().FROM(rental.AS())
	if query.PP() != "SELECT * FROM rental r" {
		t.Fail()
	}
}

func TestCreateTable(t *testing.T) {
}