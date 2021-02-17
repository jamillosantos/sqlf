package sqlf_test

import (
	"errors"
	"strings"

	"github.com/jamillosantos/sqlf"
	"github.com/jamillosantos/sqlf/testingutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operators", func() {
	Describe("And", func() {
		It("should generate an AND clause", func() {
			sb, args := new(strings.Builder), make([]interface{}, 0)
			err := sqlf.And(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()")).ToSQLFast(sb, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sb.String()).To(Equal("(id = ? AND expire_at < now())"))
		})

		It("should fail generating an errored AND clause", func() {
			sb, args := new(strings.Builder), make([]interface{}, 0)
			err := sqlf.And(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).ToSQLFast(sb, &args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})

	Describe("Or", func() {
		It("should generate an OR clause", func() {
			sb, args := new(strings.Builder), make([]interface{}, 0)
			err := sqlf.Or(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()")).ToSQLFast(sb, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sb.String()).To(Equal("(id = ? OR expire_at < now())"))
		})

		It("should fail generating an errored OR clause", func() {
			sb, args := new(strings.Builder), make([]interface{}, 0)
			err := sqlf.Or(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).ToSQLFast(sb, &args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})

	Describe("Not", func() {
		It("should generate an Not clause", func() {
			sb, args := new(strings.Builder), make([]interface{}, 0)
			err := sqlf.Not(sqlf.And(sqlf.Condition("id = ?", 1), sqlf.Condition("expire_at < now()"))).ToSQLFast(sb, &args)
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sb.String()).To(Equal("NOT (id = ? AND expire_at < now())"))
		})

		It("should fail generating an errored OR clause", func() {
			sb, args := new(strings.Builder), make([]interface{}, 0)
			err := sqlf.Not(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).ToSQLFast(sb, &args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})
})
