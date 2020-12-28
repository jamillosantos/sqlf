package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("GroupByClause", func() {
	It("should generate a GROUP BY clause", func() {
		gb := new(sqlf.GroupByClause)
		sql, args, err := gb.Fields("city", "state").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal(" GROUP BY city, state"))
	})

	It("should fail generating a GROUP BY clause with a errored field", func() {
		gb := new(sqlf.GroupByClause)
		sql, args, err := gb.Fields(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}, "state").ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a GROUP BY clause with HAVING", func() {
		gb := new(sqlf.GroupByClause).Fields("city", "state")
		gb.Having("age >= ?", 18)
		sql, args, err := gb.ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(18))
		Expect(sql).To(Equal(" GROUP BY city, state HAVING age >= ?"))
	})

	It("should generate a GROUP BY clause with HAVING and multiple criteria", func() {
		gb := new(sqlf.GroupByClause).Fields("city", "state")
		gb.HavingClause(
			sqlf.Condition("age >= ?", 18),
			sqlf.Condition("age <= ?", 35),
		)
		sql, args, err := gb.ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(18, 35))
		Expect(sql).To(Equal(" GROUP BY city, state HAVING age >= ? AND age <= ?"))
	})

	It("should fail generating a GROUP BY clause with a errored field", func() {
		gb := new(sqlf.GroupByClause).Fields("city", "state")
		gb.HavingClause(
			sqlf.Condition("age >= ?", 18),
			&testingutils.MockerSqlizer{
				SQL:  "age <= ?",
				Args: []interface{}{35},
				Err:  errors.New("forced error"),
			},
		)
		sql, args, err := gb.ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})
})
