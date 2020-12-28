package testingutils

import "strings"

// MockStringer is a mock `@fmt.Stringer` implementation that returns its `Value`
// property.
type MockStringer struct {
	Value string
}

func (s *MockStringer) String() string {
	return s.Value
}

type MockerSqlizer struct {
	SQL  string
	Args []interface{}
	Err  error
}

// ToSQL generates the SQL and returns it, alongside its params.
func (s *MockerSqlizer) ToSQL() (string, []interface{}, error) {
	return s.SQL, s.Args, s.Err
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (s *MockerSqlizer) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	if s.Err != nil {
		return s.Err
	}
	sb.WriteString(s.SQL)
	if len(s.Args) > 0 {
		*args = append(*args, s.Args...)
	}
	return nil
}
