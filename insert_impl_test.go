package sqlf_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/setare/sqlf"
	"github.com/setare/sqlf/testingutils"
)

var _ = Describe("Insert", func() {
	It("should generate a single INSERT INTO", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users", "name", "email").Values("Name 1", "email1@email.com").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"Name 1", "email1@email.com"}))
		Expect(sql).To(Equal("INSERT INTO users (name, email) VALUES (?,?)"))
	})

	It("should generate a single INSERT INTO with RETURNING", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.
			Into("users", "name", "email").
			Values("Name 1", "email1@email.com").
			Returning("id", "name").
			ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"Name 1", "email1@email.com"}))
		Expect(sql).To(Equal("INSERT INTO users (name, email) VALUES (?,?) RETURNING id, name"))
	})

	It("should fail generating a single INSERT INTO with errored field at the RETURNING clause", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.
			Into("users", "name", "email").
			Values("Name 1", "email1@email.com").
			Returning(&testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).
			ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should fail inserting an errored field", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users", "name", &testingutils.MockerSqlizer{
			SQL: "email",
			Err: errors.New("forced error"),
		}).Values("Name 1", "email1@email.com").ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a multi INSERT INTO", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users").Fields("name", "email").Values("Name 1", "email1@email.com", "Name 2", "email2@email.com").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"Name 1", "email1@email.com", "Name 2", "email2@email.com"}))
		Expect(sql).To(Equal("INSERT INTO users (name, email) VALUES (?,?), (?,?)"))
	})

	It("should fail generating a INSERT that fields and values does not match", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users", "name", "email").Values("Name 1", "email1@email.com", "Name 2").ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err).To(Equal(sqlf.ErrMismatchFieldsAndValuesCount))
	})

	It("should generate a multi INSERT INTO adding fields", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users").AddFields("name").AddFields("email").Values("Name 1", "email1@email.com").Values("Name 2", "email2@email.com").ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"Name 1", "email1@email.com", "Name 2", "email2@email.com"}))
		Expect(sql).To(Equal("INSERT INTO users (name, email) VALUES (?,?), (?,?)"))
	})

	It("should generate a INSERT INTO ... SELECT", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users", "name", "email").Select(func(s sqlf.Select) {
			s.Select("name", "email").From("employees")
		}).ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
		Expect(sql).To(Equal("INSERT INTO users (name, email) SELECT name, email FROM employees"))
	})

	It("should fail genearing a INSERT INTO ... SELECT with an error on the SELECT", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users", "name", "email").Select(func(s sqlf.Select) {
			s.Select("name", &testingutils.MockerSqlizer{
				Err: errors.New("forced error"),
			}).From("employees")
		}).ToSQL()
		Expect(err).To(HaveOccurred())
		Expect(args).To(BeNil())
		Expect(sql).To(BeEmpty())
		Expect(err.Error()).To(Equal("forced error"))
	})

	It("should generate a single INSERT INTO with suffix", func() {
		insert := new(sqlf.InsertStatement)
		sql, args, err := insert.Into("users", "name", "email").Values("Name 1", "email1@email.com").Suffix("ON DUPLICATE KEY UPDATE tries = tries + ?", 1).ToSQL()
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(Equal([]interface{}{"Name 1", "email1@email.com", 1}))
		Expect(sql).To(Equal("INSERT INTO users (name, email) VALUES (?,?) ON DUPLICATE KEY UPDATE tries = tries + ?"))
	})
})
