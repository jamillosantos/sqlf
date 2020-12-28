package sqlf

import (
	"fmt"
	"strings"
)

// RenderInterfaceAsSQL renders the input element into the `sb`(`strings.Builder`)
// according with its type, considering the input was a SQL.
//
// Sqlizer types are welcome and, if args are present they will be appended to the
// given `args` pointer.
func RenderInterfaceAsSQL(sb *strings.Builder, args *[]interface{}, element interface{}) error {
	switch p := element.(type) {
	case string:
		sb.WriteString(p)
	case fmt.Stringer:
		sb.WriteString(p.String())
	case []byte:
		sb.Write(p)
	case Sqlizer:
		err := p.ToSQLFast(sb, args)
		if err != nil {
			return err
		}
	default:
		sb.WriteString(fmt.Sprint(p))
	}
	return nil
}

// RenderInterfaceAsArg renders the input element into the `sb`(`strings.Builder`)
// according with its type considering the input as an argument.
func RenderInterfaceAsArg(sb *strings.Builder, args *[]interface{}, element interface{}) error {
	switch p := element.(type) {
	case []byte:
		sb.Write(p)
	case Sqlizer:
		err := p.ToSQLFast(sb, args)
		if err != nil {
			return err
		}
	case fmt.Stringer:
		sb.WriteString("?")
		*args = append(*args, p.String())
	case string:
		sb.WriteString("?")
		*args = append(*args, p)
	default:
		sb.WriteString("?")
		*args = append(*args, element)
	}
	return nil
}
