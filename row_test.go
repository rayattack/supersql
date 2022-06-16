package supersql_test

import "testing"

func TestRowColumn(t *testing.T) {
	res, _ := Xql.SELECT("title").FROM("film").WHERE("film_id = ?", 133).GO()
	row := res.Rows(1)

	title := row.Column("title")
	if val, ok := title.(string); !ok {
		t.Fail()
	} else {
		if val != "Chamber Italian" {
			t.Fail()
		}
	}
}
