package supersql_test

import (
	"testing"

	"github.com/rayattack/supersql"
)

var (
	actor                    = supersql.Table("actor")
	rental supersql.Relation = supersql.Table("rental")
)

func ReadRepository(table supersql.Relation) func(id string) (supersql.Results, error) {
	return func(id string) (supersql.Results, error) {
		q := Xql.SELECT().FROM(table).WHERE("actor_id = ?", id)
		results, err := q.GO()
		if err != nil {
			return nil, err
		}
		return results, nil
	}
}

func TestTableHelper(t *testing.T) {
	query := Xql.SELECT().FROM(actor).WHERE("actor_id = ?", 34)
	if query.PP() != "SELECT * FROM actor WHERE actor_id = 34" {
		t.Fail()
	}
	_, err := query.GO()
	if err != nil {
		t.Fail()
	}

	rental.AS("r")
	query = Xql.SELECT().FROM(rental.AS())
	if query.PP() != "SELECT * FROM rental r" {
		t.Fail()
	}
}

func TestCreateTable(t *testing.T) {
	select_by_id := ReadRepository(actor)
	results, err := select_by_id("45")
	if err != nil {
		t.Logf("error: %s", err)
		t.Fail()
	}
	t.Log(results.Rows(1))
}
