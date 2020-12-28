package sqlf_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/setare/sqlf"
)

var _ = Describe("Placeholder", func() {
	Describe("Dollar", func() {
		It("should replace arguments by dollar arguments", func() {
			Expect(sqlf.Dollar.Replace("SELECT * FROM users WHERE account_id = ? AND name LIKE ?")).To(Equal("SELECT * FROM users WHERE account_id = $1 AND name LIKE $2"))
		})

		It("should not replace anything", func() {
			Expect(sqlf.Dollar.Replace("SELECT * FROM users WHERE id = 4")).To(Equal("SELECT * FROM users WHERE id = 4"))
		})

		It("should escape ??", func() {
			Expect(sqlf.Dollar.Replace("SELECT * FROM users WHERE account_id = ?? AND name LIKE ??")).To(Equal("SELECT * FROM users WHERE account_id = ? AND name LIKE ?"))
		})
	})
})
