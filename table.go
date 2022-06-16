package supersql

import "fmt"

func getnom(t *SqlTable) string {
	if len(t.alias) > 0 {
		return fmt.Sprintf("%s.%s", t.alias, t.name)
	} else {
		return t.name
	}
}

type SqlTable struct {
	alias string
	name  string
	ddl   string
}
type SqlField struct{ ddl string }

func Field() Column {
	return &SqlField{``}
}

func Table(name string) (table *SqlTable) {
	table = new(SqlTable)
	table.name = name
	return
}

func (f *SqlField) IS(op string, val interface{}) string {
	return ""
}

func (t *SqlTable) AS(alias ...string) string {
	if alias != nil {
		t.alias = alias[0]
	}
	if t.alias != "" {
		return fmt.Sprintf("%s %s", t.name, t.alias)
	} else {
		return t.name
	}
}

func (t *SqlTable) DDL(ddl string) error {
	//todo implement sql syntax check logic and return early malformed warning here
	t.ddl = ddl
	return nil
}

func (t *SqlTable) GO() string {
	return t.ddl
}
