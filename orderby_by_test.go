package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("OrderByClause", func() {
	It("should generate a ORDER BY clause", func() {
		gb := new(sqlf.OrderByClause)
		sql, args, err := gb.Asc("city").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" ORDER BY city"))
	})

	It("should generate a ORDER BY clause with arguments", func() {
		gb := new(sqlf.OrderByClause)
		sql, args, err := gb.Asc(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
		}).ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(1))
		Expect(sql).To(Equal(" ORDER BY city"))
	})

	It("should fail generating a ORDER BY clause with erroed field", func() {
		gb := new(sqlf.OrderByClause)
		sql, args, err := gb.Asc(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}).ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a ORDER BY clause descending order", func() {
		gb := new(sqlf.OrderByClause)
		sql, args, err := gb.Desc("city").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" ORDER BY city DESC"))
	})

	It("should fail generating a ORDER BY clause on descending order with erroed field", func() {
		gb := new(sqlf.OrderByClause)
		sql, args, err := gb.Desc(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}).ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a ORDER BY clause with multiple fields", func() {
		gb := new(sqlf.OrderByClause)
		sql, args, err := gb.Asc("city").Desc("state").Asc("age").Desc("name").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" ORDER BY city, state DESC, age, name DESC"))
	})
})
