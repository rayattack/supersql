package supersql

import (
	"encoding/json"
	"fmt"
)

const TIP = "could not coerece data type to"

type SqlRow struct {
	datum map[string]interface{}
}

func (s SqlRow) Column(col string) interface{} {
	return s.datum[col]
}

func (s SqlRow) Scan(dest ...interface{}) error {
	return nil
}

func (s SqlRow) String(col string) (string, error) {
	val, ok := s.datum[col].(string)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Integer(col string) (int, error) {
	val, ok := s.datum[col].(int)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Boolean(col string) (bool, error) {
	val, ok := s.datum[col].(bool)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Float(col string) (float64, error) {
	val, ok := s.datum[col].(float64)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Map(col string) (map[string]interface{}, error) {
	val, ok := s.datum[col].(map[string]interface{})
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) List(col string) ([]interface{}, error) {
	val, ok := s.datum[col].([]interface{})
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Transfer(v interface{}) error {
	datum, _ := json.Marshal(s.datum)
	return json.Unmarshal(datum, v)
}

func populateRow(cols []string, values []interface{}) SqlRow {
	datum := make(map[string]interface{})
	for i, column := range cols {
		datum[column] = values[i]
	}

	return SqlRow{
		datum: datum,
	}
}
