package sqlf_test

import (
	"errors"

	"github.com/jamillosantos/sqlf"
	"github.com/jamillosantos/sqlf/testingutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Select", func() {
	Describe("Fields + From", func() {
		It("should generate select all fields", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT * FROM users"))
		})

		It("should generate select all fields from a table with alias", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users", "u").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT * FROM users AS u"))
		})

		It("should generate select all fields from a table with no alias", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users", "").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT * FROM users"))
		})

		It("should generate select all fields from a table with alias using `As` method", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").As("u").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT * FROM users AS u"))
		})

		It("should generate select some fields", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.Select("name", "email").From("users").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT name, email FROM users"))
		})

		It("should generate start adding fields", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.AddSelect("name", "email").AddSelect("age").From("users").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT name, email, age FROM users"))
		})

		It("should generate by replacing fields", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.AddSelect("name", "email").Select("age").From("users").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT age FROM users"))
		})

		It("should generate select adding fields", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.Select("name", "email").AddSelect("age").From("users").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT name, email, age FROM users"))
		})

		It("should fail generating sqlizer fields with errors", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.Select(&testingutils.MockerSqlizer{
				SQL:  "field1",
				Args: []interface{}{1},
				Err:  errors.New("forced error"),
			}).From("users").ToSQL()
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})

	Describe("Distinct", func() {
		It("should generate a distinct select", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.Select("city").Distinct().From("users").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT DISTINCT city FROM users"))
		})
	})

	Describe("Joins", func() {
		It("should generate a SQL with a JOIN", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.Select("u.*").From("users", "u").JoinClause("INNER", "permissions AS p").On("p.user_id = u.id").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u INNER JOIN permissions AS p ON p.user_id = u.id"))
		})

		It("should generate a SQL with multiple JOINs", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				JoinClause("INNER", "roles", "r").On("r.id = u.role_id").
				JoinClause("INNER", "permissions", "p").On("p.role_id = r.id AND r.type = ?", "default").
				ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(ConsistOf("default"))
			Expect(sql).To(Equal("SELECT u.* FROM users AS u INNER JOIN roles AS r ON r.id = u.role_id INNER JOIN permissions AS p ON p.role_id = r.id AND r.type = ?"))
		})

		It("should fail generating a SQL with a JOIN errored condition", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				JoinClause("INNER", "permissions", "p").OnClause(&testingutils.MockerSqlizer{
				SQL:  "p.user_id = u.id",
				Args: []interface{}{1},
				Err:  errors.New("forced error"),
			}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})

		It("should generate a INNER JOIN", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				InnerJoin("permissions", "p").On("p.user_id = u.id").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u INNER JOIN permissions AS p ON p.user_id = u.id"))
		})

		It("should generate a OUTER JOIN", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				OuterJoin("permissions", "p").On("p.user_id = u.id").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u OUTER JOIN permissions AS p ON p.user_id = u.id"))
		})

		It("should generate a LEFT JOIN", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				LeftJoin("permissions", "p").On("p.user_id = u.id").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u LEFT JOIN permissions AS p ON p.user_id = u.id"))
		})

		It("should generate a RIGHT JOIN", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				RightJoin("permissions", "p").On("p.user_id = u.id").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u RIGHT JOIN permissions AS p ON p.user_id = u.id"))
		})
	})

	Describe("Where", func() {
		It("should generate a simple where", func() {
			sql, args, err := new(sqlf.SelectStatement).
				From("users").
				Where("id = ?", 1).
				Where("age >= ?", 18).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1, 18}))
			Expect(sql).To(Equal("SELECT * FROM users WHERE id = ? AND age >= ?"))
		})

		It("should fail generating a WHERE with erroed conditions", func() {
			sql, args, err := new(sqlf.SelectStatement).
				From("users").
				WhereCriteria(&testingutils.MockerSqlizer{
					SQL:  "id = ?",
					Args: []interface{}{1},
					Err:  errors.New("forced error"),
				}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})
	})

	Describe("Group By", func() {
		It("should generate with GROUP BY", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				GroupBy("city").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u GROUP BY city"))
		})

		It("should reset GROUP BY", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				GroupBy("city").GroupBy().ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u"))
		})

		It("should generate with GROUP BY with args", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				GroupBy("city", &testingutils.MockerSqlizer{
					SQL:  "state",
					Args: []interface{}{1},
				}).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(ConsistOf(1))
			Expect(sql).To(Equal("SELECT u.* FROM users AS u GROUP BY city, state"))
		})

		It("should generate with GROUP BY with erroed fields", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.
				Select("u.*").
				From("users").As("u").
				GroupBy("city", &testingutils.MockerSqlizer{
					SQL:  "state",
					Args: []interface{}{1},
					Err:  errors.New("forced error"),
				}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})

		It("should reset GROUP BY X", func() {
			s := new(sqlf.SelectStatement).
				Select("u.*").
				From("users").As("u").
				GroupByX(func(gb sqlf.GroupBy) {
					gb.
						Fields("city").
						Having("age >= ?", 18)
				})
			sql, args, err := s.ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(ConsistOf(18))
			Expect(sql).To(Equal("SELECT u.* FROM users AS u GROUP BY city HAVING age >= ?"))
		})
	})

	Describe("Order By", func() {
		It("should generate with a ORDER BY", func() {
			sql, args, err := new(sqlf.SelectStatement).
				Select("u.*").
				From("users").As("u").
				OrderBy("city").ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u ORDER BY city"))
		})

		It("should generate with a ORDER BY with arguments", func() {
			sql, args, err := new(sqlf.SelectStatement).
				Select("u.*").
				From("users").As("u").
				OrderBy(&testingutils.MockerSqlizer{
					SQL:  "city",
					Args: []interface{}{1},
				}).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(ConsistOf(1))
			Expect(sql).To(Equal("SELECT u.* FROM users AS u ORDER BY city"))
		})

		It("should generate with a ORDER BY with errored fields", func() {
			sql, args, err := new(sqlf.SelectStatement).
				Select("u.*").
				From("users").As("u").
				OrderBy(&testingutils.MockerSqlizer{
					SQL:  "city",
					Args: []interface{}{1},
					Err:  errors.New("forced error"),
				}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})

		It("should generate with a ORDER BY desc", func() {
			s := new(sqlf.SelectStatement).
				Select("u.*").
				From("users").As("u").
				OrderByX(func(orderBy sqlf.OrderBy) {
					orderBy.Desc("city")
				})
			sql, args, err := s.ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(BeEmpty())
			Expect(sql).To(Equal("SELECT u.* FROM users AS u ORDER BY city DESC"))
		})
	})

	Describe("Offset + Limit", func() {
		It("should generate with a LIMIT clause", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").Limit(10).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(ConsistOf(10))
			Expect(sql).To(Equal("SELECT * FROM users LIMIT ?"))
		})

		It("should fail generating with a errored LIMIT clause", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").Limit(&testingutils.MockerSqlizer{
				SQL:  "?",
				Args: []interface{}{10},
				Err:  errors.New("forced error"),
			}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})

		It("should generate with a LIMIT and OFFSET clause", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").Limit(10, 20).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{20, 10}))
			Expect(sql).To(Equal("SELECT * FROM users LIMIT ? OFFSET ?"))
		})

		It("should generate with an OFFSET clause", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").Offset(20).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{20}))
			Expect(sql).To(Equal("SELECT * FROM users OFFSET ?"))
		})

		It("should fail generating with a errored OFFSET clause", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").Offset(&testingutils.MockerSqlizer{
				SQL:  "?",
				Args: []interface{}{10},
				Err:  errors.New("forced error"),
			}).ToSQL()
			Expect(err).To(HaveOccurred())
			Expect(args).To(BeNil())
			Expect(sql).To(BeEmpty())
			Expect(err.Error()).To(Equal("forced error"))
		})

		It("should generate with a LIMIT and OFFSET clause `Limit` and `Offset`", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.From("users").Limit(&testingutils.MockerSqlizer{
				SQL:  "?",
				Args: []interface{}{10},
			}).Offset(&testingutils.MockerSqlizer{
				SQL:  "? + @number",
				Args: []interface{}{20},
			}).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{10, 20}))
			Expect(sql).To(Equal("SELECT * FROM users LIMIT ? OFFSET ? + @number"))
		})
	})

	Describe("Placeholder", func() {
		It("should generate a SELECT with sequential placeholders", func() {
			s := new(sqlf.SelectStatement)
			sql, args, err := s.Placeholder(sqlf.Dollar).From("users").Where("id = ?", 1).ToSQL()
			Expect(err).ToNot(HaveOccurred())
			Expect(args).To(Equal([]interface{}{1}))
			Expect(sql).To(Equal("SELECT * FROM users WHERE id = $1"))
		})
	})
})
