package sqlf_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("JoinClause", func() {
	It("should generate a JOIN with a table name", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		err := join.Table("users").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" JOIN users"))
	})

	It("should generate a JOIN with a table name alias", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		err := join.Table("users", "u").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" JOIN users AS u"))
	})

	It("should generate a JOIN with a table name alias using the `As` method", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		err := join.Table("users").As("u").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" JOIN users AS u"))
	})

	It("should generate a JOIN with a type", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		err := join.Type("INNER").Table("users").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal("INNER JOIN users"))
	})

	It("should generate a JOIN SQL with ON clause", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		join.Table("users").On("u.id = user_id")
		err := join.ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" JOIN users ON u.id = user_id"))
	})

	It("should generate a JOIN SQL with ON clause with multiple conditions", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		join.Table("users").OnClause(sqlf.Condition("u.id = user_id"), sqlf.Condition("u.role = ?", "admin"))
		err := join.ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf("admin"))
		Expect(sb.String()).To(Equal(" JOIN users ON u.id = user_id AND u.role = ?"))
	})

	It("should generate a JOIN SQL with ON clause with an errored Sqlizer", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		join.Table("users").OnClause(&testingutils.MockerSqlizer{
			SQL:  "sqlizer1 = ?",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		})
		err := join.ToSQLFast(sb, &args)
		Expect(args).To(BeEmpty())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a JOIN SQL with USING clause", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		join.Table("users").Using("field1", sqlf.Condition("condition1"), []byte("bytes1"), &testingutils.MockStringer{Value: "stringer1"}, &testingutils.MockerSqlizer{SQL: "sqlizer1", Args: []interface{}{1}})
		err := join.ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(1))
		Expect(sb.String()).To(Equal(" JOIN users USING (field1, condition1, bytes1, stringer1, sqlizer1)"))
	})

	It("should generate a JOIN SQL with USING clause with an errored Sqlizer", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		join := new(sqlf.JoinClause)
		join.Table("users").Using(&testingutils.MockerSqlizer{SQL: "sqlizer1", Args: []interface{}{1}, Err: errors.New("forced error")})
		err := join.ToSQLFast(sb, &args)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})
})
