package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("Delete", func() {
	It("should generate a DELETE", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.From("users").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal("DELETE FROM users"))
	})

	It("should generate a DELETE CASCADE", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.Cascade().From("users").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal("DELETE CASCADE FROM users"))
	})

	It("should generate a DELETE with where", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.From("users").Where("id = ?", 1).ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{1}))
		Expect(sql).To(Equal("DELETE FROM users WHERE id = ?"))
	})

	It("should generate a DELETE with multiple where", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.From("users").Where("id = ?", 1).Where("expired_at < now()").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{1}))
		Expect(sql).To(Equal("DELETE FROM users WHERE id = ? AND expired_at < now()"))
	})

	It("should generate a DELETE with where with `Sqlizer`", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.
			From("users").
			WhereClause(sqlf.Condition("id = ?", 1)).
			WhereClause(sqlf.Condition("expired_at < now()")).ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{1}))
		Expect(sql).To(Equal("DELETE FROM users WHERE id = ? AND expired_at < now()"))
	})

	It("should fail generating a DELETE with errored where", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.From("users").WhereClause(&testingutils.MockerSqlizer{
			SQL:  "id = ?",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}).ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a DELETE with suffix", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.From("users").Suffix("SUFFIX").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal("DELETE FROM users SUFFIX"))
	})

	It("should generate a DELETE with placeholders", func() {
		d := new(sqlf.DeleteStatement)
		sql, args, err := d.Placeholder(sqlf.Dollar).From("users").Where("id = ?", 1).ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{1}))
		Expect(sql).To(Equal("DELETE FROM users WHERE id = $1"))
	})
})
