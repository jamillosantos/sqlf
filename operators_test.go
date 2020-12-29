package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("Operators", func() {
	Describe("And", func() {
		It("should generate an AND clause", func() {
			sql, args, err := sqlf.And(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()")).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sql).To(Equal("(id = ? AND expire_at < now())"))
		})

		It("should fail generating an errored AND clause", func() {
			sql, args, err := sqlf.And(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})

	Describe("Or", func() {
		It("should generate an OR clause", func() {
			sql, args, err := sqlf.Or(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()")).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sql).To(Equal("(id = ? OR expire_at < now())"))
		})

		It("should fail generating an errored OR clause", func() {
			sql, args, err := sqlf.Or(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})

	Describe("Not", func() {
		It("should generate an Not clause", func() {
			sql, args, err := sqlf.Not(sqlf.And(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()"))).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sql).To(Equal("NOT (id = ? AND expire_at < now())"))
		})

		It("should fail generating an errored OR clause", func() {
			sql, args, err := sqlf.Not(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})
})
