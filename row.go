package supersql

import (
	"encoding/json"
	"fmt"
)

const TIP = "could not coerece data type to"

type SqlRow struct {
	data map[string]interface{}
}

func (s SqlRow) Column(col string) interface{} {
	return s.data[col]
}

func (s SqlRow) Scan(dest ...interface{}) error {
	return nil
}

func (s SqlRow) String(col string) (string, error) {
	val, ok := s.data[col].(string)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Integer(col string) (int, error) {
	val, ok := s.data[col].(int)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Boolean(col string) (bool, error) {
	val, ok := s.data[col].(bool)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Float(col string) (float64, error) {
	val, ok := s.data[col].(float64)
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Map(col string) (map[string]interface{}, error) {
	val, ok := s.data[col].(map[string]interface{})
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) List(col string) ([]interface{}, error) {
	val, ok := s.data[col].([]interface{})
	if ok {
		return val, nil
	}
	return val, fmt.Errorf("%s %T", TIP, val)
}

func (s SqlRow) Transfer(v interface{}) error {
	data, _ := json.Marshal(s.data)
	return json.Unmarshal(data, v)
}

func populateRow(cols []string, values []interface{}) SqlRow {
	data := make(map[string]interface{})
	for i, column := range cols {
		data[column] = values[i]
	}

	return SqlRow{
		data: data,
	}
}
