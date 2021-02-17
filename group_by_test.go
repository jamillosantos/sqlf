package sqlf_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamillosantos/sqlf"
	"github.com/jamillosantos/sqlf/testingutils"
)

var _ = Describe("GroupByClause", func() {
	It("should generate a GROUP BY clause", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.GroupByClause)
		err := gb.Fields("city", "state").ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sb.String()).To(Equal(" GROUP BY city, state"))
	})

	It("should fail generating a GROUP BY clause with a errored field", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.GroupByClause)
		err := gb.Fields(&testingutils.MockerSqlizer{
			SQL:  "city",
			Args: []interface{}{1},
			Err:  errors.New("forced error"),
		}, "state").ToSQLFast(sb, &args)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a GROUP BY clause with HAVING", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.GroupByClause).Fields("city", "state")
		gb.Having("age >= ?", 18)
		err := gb.ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(18))
		Expect(sb.String()).To(Equal(" GROUP BY city, state HAVING age >= ?"))
	})

	It("should generate a GROUP BY clause with HAVING and multiple criteria", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.GroupByClause).Fields("city", "state")
		gb.HavingClause(
			sqlf.Condition("age >= ?", 18),
			sqlf.Condition("age <= ?", 35),
		)
		err := gb.ToSQLFast(sb, &args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(ConsistOf(18, 35))
		Expect(sb.String()).To(Equal(" GROUP BY city, state HAVING age >= ? AND age <= ?"))
	})

	It("should fail generating a GROUP BY clause with a errored field", func() {
		sb, args := new(strings.Builder), make([]interface{}, 0)
		gb := new(sqlf.GroupByClause).Fields("city", "state")
		gb.HavingClause(
			sqlf.Condition("age >= ?", 18),
			&testingutils.MockerSqlizer{
				SQL:  "age <= ?",
				Args: []interface{}{35},
				Err:  errors.New("forced error"),
			},
		)
		err := gb.ToSQLFast(sb, &args)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("forced error"))
	})
})
