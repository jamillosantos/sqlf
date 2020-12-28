package sqlf

import (
	"strconv"
	"strings"
)

// PlaceholderFormat enables the system to use different placeholder formats.
type PlaceholderFormat interface {
	Replace(sql string) (string, error)
}

type questionPlaceholder struct{}

func (*questionPlaceholder) Replace(sql string) (string, error) {
	return sql, nil
}

type dollarPlaceholder struct{}

var Dollar = &dollarPlaceholder{}

var questionSign = byte('?')

func (*dollarPlaceholder) Replace(sql string) (string, error) {
	lenSQL := len(sql)
	count := 0
	sb := new(strings.Builder)
	lastI := 0
	for i := 0; i < lenSQL; i++ {
		if sql[i] == questionSign {
			if i+1 < lenSQL && sql[i+1] == questionSign {
				sb.WriteString(sql[lastI:i])
				i++
				lastI = i
				continue
			}
			count++
			sb.WriteString(sql[lastI:i])
			sb.WriteByte('$')
			sb.WriteString(strconv.Itoa(count))
			lastI = i + 1
		}
	}
	sb.WriteString(sql[lastI:])
	return sb.String(), nil
}
