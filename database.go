package supersql

type SqlDatabase struct {
	ddl string
}

func (d *SqlDatabase) DDL(ddl string) {
	d.ddl = ddl
}

func Database() *SqlDatabase {
	return &SqlDatabase{}
}
