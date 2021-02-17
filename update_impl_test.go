package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamillosantos/sqlf"
	"github.com/jamillosantos/sqlf/testingutils"
)

var _ = Describe("Update", func() {
	It("should generate a UPDATE", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.Table("users").Set("name", "name1").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1"}))
		Expect(sql).To(Equal("UPDATE users SET name = ?"))
	})

	It("should generate a UPDATE with no alias", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.Table("users", "").Set("name", "name1").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1"}))
		Expect(sql).To(Equal("UPDATE users SET name = ?"))
	})

	It("should generate a UPDATE with alias", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.Table("users", "u").Set("u.name", "name1").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1"}))
		Expect(sql).To(Equal("UPDATE users AS u SET u.name = ?"))
	})

	It("should generate a UPDATE with multiple fields", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.
			Table("users").
			Set(
				"name", "name1",
				"email", "email1@email.com",
			).
			ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1", "email1@email.com"}))
		Expect(sql).To(Equal("UPDATE users SET name = ?, email = ?"))
	})

	It("should generate a UPDATE with errored field name", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.Table("users").Set(&testingutils.MockerSqlizer{
			Err: errors.New("forced error"),
		}, "name1").ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a UPDATE with errored value", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.Table("users").Set(
			"name",
			&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).
			ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should fail generating an UPDATE with wrong field and values count", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.
			Table("users").
			Set(
				"name", "name1",
				"email",
			).
			ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err).To(Equal(sqlf.ErrUpdateInvalidFieldValuePairCount))
	})

	It("should generate a UPDATE with the WHERE clause", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.
			Table("users").
			Set("name", "name1").
			Set("email", "email1@email.com").
			Where("id = ?", 1).
			ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1", "email1@email.com", 1}))
		Expect(sql).To(Equal("UPDATE users SET name = ?, email = ? WHERE id = ?"))
	})

	It("should fail generating a UPDATE with a errored WHERE clause", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.
			Table("users").
			Set("name", "name1").
			Set("email", "email1@email.com").
			WhereClause(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).
			ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a UPDATE with the multiple WHERE conditions", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.
			Table("users").
			Set("name", "name1").
			Set("email", "email1@email.com").
			Where("id = ?", 1).
			Where("expire_at > now()").
			ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1", "email1@email.com", 1}))
		Expect(sql).To(Equal("UPDATE users SET name = ?, email = ? WHERE id = ? AND expire_at > now()"))
	})

	It("should generate a UPDATE with the multiple WHERE clauses", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.
			Table("users").
			Set("name", "name1").
			Set("email", "email1@email.com").
			WhereClause(sqlf.Or(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()"))).
			WhereClause(sqlf.Condition("enabled = ?", 1)).
			ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1", "email1@email.com", 1, 1}))
		Expect(sql).To(Equal("UPDATE users SET name = ?, email = ? WHERE (id = ? OR expire_at < now()) AND enabled = ?"))
	})

	It("should generate a UPDATE with a placeholder", func() {
		d := new(sqlf.UpdateStatement)
		sql, args, err := d.Placeholder(sqlf.Dollar).Table("users").Set("name", "name1").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"name1"}))
		Expect(sql).To(Equal("UPDATE users SET name = $1"))
	})

})
