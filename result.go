package supersql

type SqlResult struct {
	cols  []string
	rows  []Row
	count int
}

func (r SqlResult) All() []Row {
	return r.rows
}

func (r SqlResult) Count() int {
	return r.count
}

//Rows allow for querying returned results as it would appear in a database view style i.e. the first
//result being at position 1 and not 0, and so on and so forth.
//This method will always return the last value if the provided position is greater than the results fetched
//from the database.
func (r SqlResult) Rows(position int) Row {
	if position > r.Count() {
		position = r.Count()
	}
	return r.rows[position - 1]
}

//TODO documentation for Transfer
func (r SqlResult) Transfer(s []map[string]interface{}) error {
	return nil
}

// TODO
// Unlike SqlResults this does not load all results in to memory but works with
// the Next() paradigm of package sql etc.
type SqlStream struct {
}
