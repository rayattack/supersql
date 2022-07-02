package supersql

import "fmt"

type Field struct {
	name   string
	ddl    string
	quoted bool
}

func colmaker(name string, options ...interface{}) string {
	return ""
}

func operator(f Field, op string, val interface{}) string {
	var eval string
	if f.quoted {
		eval = fmt.Sprintf("'%v'", val)
	} else {
		eval = fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("%s %s %v", f.name, op, eval)
}

func (f Field) Eq(val interface{}) string {
	return operator(f, "=", val)
}

func Integer(name string, options ...interface{}) Field {
	return Field{name, colmaker("integer"), false}
}

func Varchar(name string) Field {
	return Field{name, fmt.Sprintf("%s  varchar", name), true}
}
