package supersql

type Relation interface {
	AS(alias ...string) string
	DDL(ddl string) error
}

type Command interface {
	// AS(alias string) Command
	// ASC(col string) Command
	// CREATE(entity interface{}) Command
	// DESC(col string) Command
	FROM(entities ...interface{}) Command
	GO() (Results, error)
	// INSERT(columns ...string) Command
	// INSERT_INTO(table interface{}, columns ...string) Command
	// INTO(entity interface{}) Command
	JOIN(entity interface{}) Command
	// LIMIT(count int) Command
	// OFFSET(count int) Command
	ON(statement string, conditions ...interface{}) Command
	PP() string
	SELECT(columns ...string) Command
	// VALUES(values [][]interface{}) Command
	WHERE(statement string, conditions ...interface{}) Command
}

type Db interface {
	DDL(ddl string) error
}

type Column interface {
	IS(op string, val interface{}) string
}

type Results interface {
	All() []Row
	Count() int
	Rows(position int) Row
}

type Row interface {
	Column(col string) interface{}
	Scan(dest ...interface{}) error
	String(col string) (string, error)
	Integer(col string) (int, error)
	Boolean(col string) (bool, error)
	Float(col string) (float64, error)
	Map(col string) (map[string]interface{}, error)
	List(col string) ([]interface{}, error)
	Transfer(v interface{}) error
}

type Values [][]interface{}