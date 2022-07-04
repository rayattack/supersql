package supersql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const POSTGRES_MAX_COLUMNS = 1600

type SqlQuery struct {
	conn *pgx.Conn
	pool *pgxpool.Pool
	ssql string
	ctx  context.Context
	args []interface{}
	cols []string
	vals [][]interface{}
	void bool
}

//TODO: AS Documentation
func (q SqlQuery) AS(alias string) Command {
	q.ssql = fmt.Sprintf("%s AS %s ", q.ssql, alias)
	return &q
}

//TODO: ASC Documentation
func (q SqlQuery) ASC(ob ...string) Command {
	if ob != nil {
		return q.ob(fmt.Sprintf("ORDER BY %s ASC", ob[0]))
	}
	q.ssql = fmt.Sprintf("%s ASC", q.ssql)
	return q
}

//TODO: CLOSE Documentation
func (q *SqlQuery) CLOSE() error {
	if q.pool != nil {
		q.pool.Close()
		return nil
	}
	return q.conn.Close(q.ctx)
}

//TODO: DESC Documentation
func (q SqlQuery) DESC(ob ...string) Command {
	if ob != nil {
		return q.ob(fmt.Sprintf("ORDER BY %s DESC", ob[0]))
	}
	q.ssql = fmt.Sprintf("%s DESC", q.ssql)
	return q
}

//Helper function for DRY purposes to optimize inserting records by using pgx.CopyFrom as opposed
//to a naive SQL INSERT command
func (q SqlQuery) do(rows [][]interface{}) error {
	affected, error := q.conn.CopyFrom(context.TODO(), pgx.Identifier{}, q.cols, pgx.CopyFromRows(rows))
	if error != nil {
		return error
	}
	fmt.Print(affected)
	return nil
}

//TODO: FROM Documentation
func (q SqlQuery) FROM(entities ...interface{}) Command {
	e := []string{}
	for _, entity := range entities {
		t := coerceToString(entity)
		e = append(e, t)
	}
	q.ssql = fmt.Sprintf("%s FROM %s", q.ssql, strings.Join(e, ","))
	return q
}

//Send your expantiated query to the server for execution. If an integer argument is provided
//then results up to the value specified will be preloaded in an sqlresult object. If the optional
//integer argument is less than zero i.e. -1, then all results will be loaded and if none is provided
//then calling next() on the sqlresult and streaming back will be the default behaviour activated.
func (q SqlQuery) GO(prefetch ...int) (Results, error) {
	//i.e. if cols present we are in insert mode
	if q.cols != nil {
		return nil, q.do(q.vals)
	}

	q.ssql = countAndReplacePlaceholders(q.ssql)

	if q.void {
		_, error := q.pool.Exec(q.ctx, q.ssql, q.args...)
		if error != nil {
			return nil, error
		}
		return nil, nil
	}
	ctrl, err := q.pool.Query(q.ctx, q.ssql, q.args...)
	if err != nil {
		return nil, err
	}

	var columns []string
	for _, column := range ctrl.FieldDescriptions() {
		columns = append(columns, string(column.Name))
	}

	//Check for streaming flag here and stream (lazy read) as opposed to greedy read
	rows := []Row{}
	count := 0
	for ctrl.Next() {
		values, err := ctrl.Values()
		if err != nil {
			return nil, err
		}
		rows = append(rows, populateRow(columns, values))
		count++
	}
	return SqlResult{columns, rows, count}, nil
}

//Only use this function if q.INTO(...) will be invoked immediately after this
//is called i.e. this will register the value passed in for inversion
//when q.INTO(...) is invoked
func (q SqlQuery) INSERT(columns ...string) Command {
	q.cols = columns
	q.ssql = fmt.Sprintf("(%s)", strings.Join(columns, ", "))
	return q
}

//Responsible for expantiation sql to write or create new records. This command is
//exactly the same as invoking q.INSERT(...) followed immediately by q.INTO(...)
func (q SqlQuery) INSERT_INTO(table interface{}, optionalColumns ...[]string) Command {
	t := coerceToString(table)

	if len(optionalColumns) > 0 {
		columns := optionalColumns[0]
		if interpolative := len(columns); interpolative > 0 {
			q.cols = columns
			if placeholders := strings.Count(t, "?"); placeholders > 0 {
				for _, col := range columns {
					t = strings.Replace(t, "?", col, 1)
				}
			} else {
				t = fmt.Sprintf("%s (%s)", t, strings.Join(columns, ", "))
			}
		} else {
			_, cols, _ := strings.Cut(t, "(")
			q.cols = strings.Split(cols[:len(cols)-1], ",")
		}
	}
	q.ssql = fmt.Sprintf("INSERT INTO %s", t)
	return q
}

//Only use this function if q.INSERT(...) invoked immediately before this
//was called i.e. at the point this is invoked: q.ssql = q.INSERT()
//and so q.INTO(...) inverts the order autocorrecting q.ssql for further
//expantiation
func (q SqlQuery) INTO(table interface{}) Command {
	t := coerceToString(table)
	//revisit and check correctness
	return q.INSERT_INTO(fmt.Sprintf("%s %s", t, q.ssql))
}

//Issue a join SQL command to tie entities/tables together. This should always
//be followed by an invokation of q.ON(...)
func (q SqlQuery) JOIN(entity interface{}) Command {
	t := coerceToString(entity)
	q.ssql = fmt.Sprintf("%s JOIN %s", q.ssql, t)
	return q
}

//TODO: LIMIT Documentation
func (q SqlQuery) LIMIT(limit int) Command {
	q.ssql = fmt.Sprintf("%s LIMIT %d", q.ssql, limit)
	return &q
}

//TODO: OFFSET Documentation
func (q SqlQuery) OFFSET(offset int) Command {
	q.ssql = fmt.Sprintf("%s OFFSET %d", q.ssql, offset)
	return q
}

//Continuation expantiator for JOIN(...) SQL command. This function provides
//a simple way to specify how the entities should be joined i.e. what columns across
//the two entities intersect
func (q SqlQuery) ON(statement string, conditions ...interface{}) Command {
	q.args = append(q.args, conditions...)
	q.ssql = fmt.Sprintf("%s ON %s", q.ssql, statement)
	return q
}

//Order by helper
func (q SqlQuery) ob(ob string) Command {
	q.ssql = fmt.Sprintf("%s %s", q.ssql, ob)
	return q
}

//TODO: ORDER_BY Documentation
func (q SqlQuery) ORDER_BY(ob string) Command {
	q.ssql = fmt.Sprintf("%s ORDER BY %s", q.ssql, ob)
	return q
}

//(PP = PrettyPrint) Returns whatever sql statement has been expantiated at the point this function
//is invoked.
func (q SqlQuery) PP() string {
	csql := q.ssql
	for _, arg := range q.args {
		csql = strings.Replace(csql, "?", fmt.Sprint(arg), 1)
	}
	return csql
}

func (q SqlQuery) RUN(ddl string) Command {
	q.ssql = ddl
	q.void = true
	return q
}

//TODO: SELECT Documentation
func (q SqlQuery) SELECT(fields ...string) Command {
	q.void = false
	if len(fields) == 0 {
		fields = []string{"*"}
	}
	q.ssql = fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
	return q
}

//An array(golang slice) contaning another array(slice) of dynamic values to be used in adding
//a new record.
//TIP: you can create an alias type and OR an expansion function helper and reuse
//that in your own code to save a few keystrokes i.e.
//In the current iteration this code is semantic only i.e. not used as it shortcircuits in
//the GO() method invokation to use a more performant
//golang psql specific sql construct. Take a look at the GO method for more details.
//This might change if more database libraries are supported in another version
func (q SqlQuery) VALUES(vals ...[]interface{}) Command {
	q.vals = vals
	values := []string{}
	for _, val := range vals {
		replacement := strings.Repeat("?, ", len(val)-1)
		replacement = fmt.Sprintf("%s?", replacement)
		values = append(values, fmt.Sprintf("(%s)", replacement))
	}
	q.ssql = fmt.Sprintf("%s VALUES %s", q.ssql, strings.Join(values, ","))
	return q
}

//TODO: WHERE Documentation
func (q SqlQuery) WHERE(statement string, conditions ...interface{}) Command {
	q.args = append(q.args, conditions...)
	q.ssql = fmt.Sprintf("%s WHERE %s", q.ssql, statement)
	return q
}

//TODO: Query Documentation
func Query(ctx context.Context, dsn string) (*SqlQuery, error) {
	var pool *pgxpool.Pool
	var conn *pgx.Conn
	var ssql string
	var args []interface{}
	var cols []string
	var vals [][]interface{}

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Could not initialize connection to database due to: %s", err)
		return nil, err
	}

	query := &SqlQuery{
		conn: conn,
		pool: pool,
		ssql: ssql,
		ctx:  ctx,
		args: args,
		cols: cols,
		vals: vals,
		void: true,
	}
	return query, nil
}
