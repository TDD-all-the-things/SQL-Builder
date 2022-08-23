package builder

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrInvalidEntity = errors.New("builder: invalid entity")
)

func InsertStmt(entity any) (query string, values []any, err error) {
	insert := NewInsertSQL()
	insert.Fill(reflect.ValueOf(entity))
	return insert.String(), insert.Values(), insert.Err()
}

type InsertSQL struct {
	name    string
	columns []string
	values  []any
	seen    map[string]struct{}
	err     error
}

func NewInsertSQL() *InsertSQL {
	return &InsertSQL{seen: make(map[string]struct{})}
}

func (t *InsertSQL) Fill(val reflect.Value) {

	if !val.IsValid() {
		t.err = fmt.Errorf("%w", ErrInvalidEntity)
		return
	}

	// multiple level pointer
	for i := 0; val.Kind() == reflect.Pointer; i++ {
		val = val.Elem()
		if i > 0 {
			t.err = fmt.Errorf("%w", ErrInvalidEntity)
			return
		}
	}

	// (*struct)(nil) or non-struct or empty struct
	if !val.IsValid() || val.Kind() != reflect.Struct || val.NumField() == 0 {
		t.err = fmt.Errorf("%w", ErrInvalidEntity)
		return
	}

	if t.name == "" {
		t.SetName(val.Type().Name())
	}

	for i := 0; i < val.NumField(); i++ {

		if t.HasColumn(val.Type().Field(i).Name) {
			continue
		}

		isEmbeddedField := val.Type().Field(i).Anonymous
		isStruct := val.Field(i).Kind() == reflect.Struct
		_, implementsInterface := val.Field(i).Interface().(driver.Valuer)

		if isEmbeddedField && isStruct && !implementsInterface {
			t.Fill(val.Field(i))
			continue
		}

		t.AddColumn(val.Type().Field(i).Name, val.Field(i).Interface())
	}
}

func (t *InsertSQL) SetName(name string) {
	t.name = "`" + name + "`"
}

func (t *InsertSQL) AddColumn(name string, val any) {
	column := "`" + name + "`"
	t.columns = append(t.columns, column)
	t.seen[column] = struct{}{}
	t.values = append(t.values, val)
}

func (t *InsertSQL) HasColumn(name string) bool {
	_, ok := t.seen["`"+name+"`"]
	return ok
}

func (t *InsertSQL) Values() []any {
	if t.err != nil {
		return nil
	}
	return t.values
}

func (t *InsertSQL) Err() error {
	return t.err
}

func (t *InsertSQL) String() string {
	if t.name == "" || t.columns == nil || t.err != nil {
		return ""
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", t.name, strings.Join(t.columns, ","),
		strings.TrimRight(strings.Repeat("?,", len(t.seen)), ","))
}
