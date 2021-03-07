package sqlf

import (
	"strconv"
	"strings"
	"sync"
)

// Placeholders returns a sequence of "?" separated by comma.
func Placeholders(count int) string {
	return strings.Repeat(",?", count)[1:]
}

// PlaceholderFormatFactory enables the system to use different placeholder formats. It wraps the SQLWriter into another
// that will replace ? by the desired placeholder.
type PlaceholderFormatFactory interface {
	// Wrap wraps the given `SQLWriter` into the placeholder replacer.
	Wrap(writer SQLWriter) SQLWriter

	// Put aims to return the wrapper to a pool. Of course, it depends on the factory implementation.
	Put(writer SQLWriter)
}

type questionPlaceholderFactory struct{}

type dollarPlaceholderFactory struct {
	pool sync.Pool
}

type dollarPlaceholder struct {
	writer           SQLWriter
	placeholderCount int
}

var (
	QuestionPlaceholder = &questionPlaceholderFactory{}
	DollarPlaceholder   = &dollarPlaceholderFactory{
		pool: sync.Pool{
			New: func() interface{} {
				return &dollarPlaceholder{}
			},
		},
	}
)

var (
	questionSign = byte('?')
	dollarSign   = byte('$')
)

// Wrap returns the given sqlWriter without wrapping it to anything. That is so because the default placeholder is
// the `?`. So, nothing need to be replaced.
func (q *questionPlaceholderFactory) Wrap(sqlWriter SQLWriter) SQLWriter {
	return sqlWriter
}

// Put just does nothing as `?` placeholder is the default.
func (q *questionPlaceholderFactory) Put(SQLWriter) {}

// Wrap wraps the given `sqlWriter` into a `dollarPlaceholder` that will replace any found `?` by a `$1` where 1 is
// the index of the placeholder. Each `?` will be considered a new placeholder.
//
// To add `?` to the SQL query, you should double `??`. This way, the placeholder will escape and output `?`.
func (q *dollarPlaceholderFactory) Wrap(sqlWriter SQLWriter) SQLWriter {
	writer := q.pool.Get().(*dollarPlaceholder)
	writer.writer = sqlWriter
	return writer
}

func (q *dollarPlaceholderFactory) Put(sqlWriter SQLWriter) {
	dollar, ok := sqlWriter.(*dollarPlaceholder)
	if !ok {
		return
	}

	dollar.placeholderCount = 0
	dollar.writer = nil
}

func (dp *dollarPlaceholder) WriteByte(p byte) (err error) {
	return dp.writer.WriteByte(p)
}

func (dp *dollarPlaceholder) Write(p []byte) (n int, err error) {
	lastW := 0
	for i := 0; i < len(p); i++ {
		isInterrogation := p[i] == '?'
		if !isInterrogation {
			continue
		}

		if isInterrogation && len(p) > i+1 && p[i+1] == '?' {
			// Write to the original writer and skip a ?
			dp.writer.Write(p[lastW:i])
			i++
			lastW = i
			continue
		}

		if isInterrogation {
			dp.placeholderCount++
			_, err = dp.writer.Write(p[lastW:i])
			if err != nil {
				return 0, err // Well, this is not true. But it will do for now.
			}
			err = dp.writer.WriteByte(dollarSign)
			if err != nil {
				return 0, err // Well, this is not true. But it will do for now.
			}
			_, err = dp.writer.WriteString(strconv.Itoa(dp.placeholderCount))
			if err != nil {
				return 0, err // Well, this is not true. But it will do for now.
			}
			lastW = i + 1
		}
	}

	// If there is something in the buffer to be written.
	if lastW != len(p) {
		dp.writer.Write(p[lastW:])
	}
	return len(p), nil
}

func (dp *dollarPlaceholder) WriteString(s string) (n int, err error) {
	lastW := 0
	for i := 0; i < len(s); i++ {
		isInterrogation := s[i] == '?'
		if !isInterrogation {
			continue
		}

		if isInterrogation && len(s) > i+1 && s[i+1] == '?' {
			// Write to the original writer and skip a ?
			dp.writer.WriteString(s[lastW:i])
			i++
			lastW = i
			continue
		}

		if isInterrogation {
			dp.placeholderCount++
			_, err = dp.writer.WriteString(s[lastW:i])
			if err != nil {
				return 0, err // Well, this is not true. But it will do for now.
			}
			err = dp.writer.WriteByte(dollarSign)
			if err != nil {
				return 0, err // Well, this is not true. But it will do for now.
			}
			_, err = dp.writer.WriteString(strconv.Itoa(dp.placeholderCount))
			if err != nil {
				return 0, err // Well, this is not true. But it will do for now.
			}
			lastW = i + 1
		}
	}

	// If there is something in the buffer to be written.
	if lastW != len(s) {
		dp.writer.WriteString(s[lastW:])
	}
	return len(s), nil
}

func (dp *dollarPlaceholder) String() string {
	return dp.writer.String()
}
