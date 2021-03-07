package sqlf_test

import (
	"strings"

	"github.com/jamillosantos/sqlf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Placeholder", func() {
	Describe("Dollar", func() {
		Describe("for string", func() {
			It("should replace arguments by dollar arguments", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.WriteString("SELECT * FROM users WHERE account_id = ? AND name LIKE ?")
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = $1 AND name LIKE $2"))
			})

			It("should replace arguments by dollar arguments in multiple parts", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.WriteString("SELECT * FROM users WHERE account_id = ?")
				Expect(err).ToNot(HaveOccurred())
				_, err = ph.WriteString(" AND name LIKE ?")
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = $1 AND name LIKE $2"))
			})

			It("should not replace anything", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.WriteString("SELECT * FROM users WHERE id = 4")
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE id = 4"))
			})

			It("should escape ?", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.WriteString("SELECT * FROM users WHERE account_id = ?? AND name LIKE ??")
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = ? AND name LIKE ?"))
			})

			It("should escape sequential ?", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.WriteString("SELECT * FROM users WHERE account_id = ?????? AND name LIKE ????")
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = ??? AND name LIKE ??"))
			})
		})

		Describe("for bytes", func() {
			It("should replace arguments by dollar arguments", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.Write([]byte("SELECT * FROM users WHERE account_id = ? AND name LIKE ?"))
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = $1 AND name LIKE $2"))
			})

			It("should replace arguments by dollar arguments in multiple parts", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.Write([]byte("SELECT * FROM users WHERE account_id = ?"))
				Expect(err).ToNot(HaveOccurred())
				_, err = ph.Write([]byte(" AND name LIKE ?"))
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = $1 AND name LIKE $2"))
			})

			It("should not replace anything", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.Write([]byte("SELECT * FROM users WHERE id = 4"))
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE id = 4"))
			})

			It("should escape ?", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.Write([]byte("SELECT * FROM users WHERE account_id = ?? AND name LIKE ??"))
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = ? AND name LIKE ?"))
			})

			It("should escape sequential ?", func() {
				sb := new(strings.Builder)
				ph := sqlf.DollarPlaceholder.Wrap(sb)
				_, err := ph.Write([]byte("SELECT * FROM users WHERE account_id = ?????? AND name LIKE ????"))
				Expect(err).ToNot(HaveOccurred())
				Expect(ph.String()).To(Equal("SELECT * FROM users WHERE account_id = ??? AND name LIKE ??"))
			})
		})
	})
})
