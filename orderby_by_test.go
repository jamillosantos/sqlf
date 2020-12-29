package sqlf_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("OrderByClause", func() {
	It("should generate a ORDER BY clause", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.OrderByClause)
		err := gb.Asc("city").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" ORDER BY city"))
	})

	It("should generate a ORDER BY clause with arguments", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.OrderByClause)
		err := gb.Asc(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
		}).ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(1))
		Expect(sb.String()).To(Equal(" ORDER BY city"))
	})

	It("should fail generating a ORDER BY clause with erroed field", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.OrderByClause)
		err := gb.Asc(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}).ToSQLFast(sb, &args)
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a ORDER BY clause descending order", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.OrderByClause)
		err := gb.Desc("city").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" ORDER BY city DESC"))
	})

	It("should fail generating a ORDER BY clause on descending order with erroed field", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.OrderByClause)
		err := gb.Desc(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}).ToSQLFast(sb, &args)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a ORDER BY clause with multiple fields", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.OrderByClause)
		err := gb.Asc("city").Desc("state").Asc("age").Desc("name").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" ORDER BY city, state DESC, age, name DESC"))
	})
})
