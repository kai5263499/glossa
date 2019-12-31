package glossa

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type matchTest struct {
	content   string
	pattern   Pattern
	callback  func(...[]interface{})
	match     bool
	expectErr bool
}

var _ = Describe("pattern", func() {
	It("Should create patterns properly", func() {
		var err error
		Expect(err).To(BeNil())
	})
})
