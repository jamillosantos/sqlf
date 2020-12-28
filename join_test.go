package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("JoinClause", func() {
	It("should generate a JOIN with a table name", func() {
		join := new(sqlf.JoinClause)
		sql, args, err := join.Table("users").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" JOIN users"))
	})

	It("should generate a JOIN with a table name alias", func() {
		join := new(sqlf.JoinClause)
		sql, args, err := join.Table("users", "u").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" JOIN users AS u"))
	})

	It("should generate a JOIN with a table name alias using the `As` method", func() {
		join := new(sqlf.JoinClause)
		sql, args, err := join.Table("users").As("u").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" JOIN users AS u"))
	})

	It("should generate a JOIN with a type", func() {
		join := new(sqlf.JoinClause)
		sql, args, err := join.Type("INNER").Table("users").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal("INNER JOIN users"))
	})

	It("should generate a JOIN SQL with ON clause", func() {
		join := new(sqlf.JoinClause)
		join.Table("users").On("u.id = user_id")
		sql, args, err := join.ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" JOIN users ON u.id = user_id"))
	})

	It("should generate a JOIN SQL with ON clause with multiple conditions", func() {
		join := new(sqlf.JoinClause)
		join.Table("users").OnClause(sqlf.Condition("u.id = user_id"), sqlf.Condition("u.role = ?", "admin"))
		sql, args, err := join.ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf("admin"))
		Expect(sql).To(Equal(" JOIN users ON u.id = user_id AND u.role = ?"))
	})

	It("should generate a JOIN SQL with ON clause with an errored Sqlizer", func() {
		join := new(sqlf.JoinClause)
		join.Table("users").OnClause(&testingutils.MockerSqlizer{
			SQL:  "sqlizer1 = ?",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		})
		sql, args, err := join.ToSQL()
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a JOIN SQL with USING clause", func() {
		join := new(sqlf.JoinClause)
		join.Table("users").Using("field1", sqlf.Condition("condition1"), []byte("bytes1"), &testingutils.MockStringer{Value: "stringer1"}, &testingutils.MockerSqlizer{SQL: "sqlizer1", Args: []interface{}{1}})
		sql, args, err := join.ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(1))
		Expect(sql).To(Equal(" JOIN users USING (field1, condition1, bytes1, stringer1, sqlizer1)"))
	})

	It("should generate a JOIN SQL with USING clause with an errored Sqlizer", func() {
		join := new(sqlf.JoinClause)
		join.Table("users").Using(&testingutils.MockerSqlizer{SQL: "sqlizer1", Args: []interface{}{1}, Err: errors.New("forced error")})
		sql, args, err := join.ToSQL()
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})
})
