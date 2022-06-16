package supersql_test

import (
	"reflect"
	"testing"

	"github.com/rayattack/supersql"
)

func TestSingleSelect(t *testing.T) {
	r, _ := Xql.SELECT("title").FROM("film").WHERE("film_id = ?", 007).GO()
	if r.Count() != 1 {
		t.Fail()
	}

	row := r.Rows(45)
	if _, ok := row.(supersql.SqlRow); !ok {
		t.Fail()
	}

	rows := r.All()
	if !reflect.DeepEqual(rows[0], row) {
		t.Fail()
	}
}
