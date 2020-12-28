package sqlf_test

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/elgris/sqrl"
	"github.com/setare/sqlf"
)

func BenchmarkSQLFSelectCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := new(sqlf.SelectStatement)
		_, _, err := s.From("users", "u").InnerJoin("permissions", "p").On("p.user_id = u.id").Where("u.age >= ?", 18).ToSQL()
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSquirrelSelectCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := sq.Select("*").From("users u").InnerJoin("permissions ON p.user_id = u.id").Where("u.age >= ?", 18).ToSql()
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSqrlSelectCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := sqrl.Select("*").From("users u").Join("permissions ON p.user_id = u.id").Where("u.age >= ?", 18).ToSql()
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSQLFSelectSQLGeneration(b *testing.B) {
	s := new(sqlf.SelectStatement)
	s.From("users", "u").InnerJoin("permissions", "p").On("p.user_id = u.id").Where("u.age >= ?", 18)
	for i := 0; i < b.N; i++ {
		_, _, err := s.ToSQL()
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSquirrelSelectSQLGeneration(b *testing.B) {
	s := sq.Select("*").From("users u").InnerJoin("permissions ON p.user_id = u.id").Where("u.age >= ?", 18)
	for i := 0; i < b.N; i++ {
		_, _, err := s.ToSql()
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSqrlSelectSQLGeneration(b *testing.B) {
	s := sqrl.Select("*").From("users u").Join("permissions ON p.user_id = u.id").Where("u.age >= ?", 18)
	for i := 0; i < b.N; i++ {
		_, _, err := s.ToSql()
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

var sqlForPlaceholder = "SELECT * FROM users WHERE account_id = ? AND name LIKE ?"

func BenchmarkSQLFDollarPlaceholder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := sqlf.Dollar.Replace(sqlForPlaceholder)
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSquirrelDollarPlaceholder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := sq.Dollar.ReplacePlaceholders(sqlForPlaceholder)
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}

func BenchmarkSqrlDollarPlaceholder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := sqrl.Dollar.ReplacePlaceholders(sqlForPlaceholder)
		if err != nil {
			fmt.Println(err)
			b.Fail()
			return
		}
	}
}
