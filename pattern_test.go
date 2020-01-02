package glossa

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type matchTest struct {
	content      string
	pattern      Pattern
	argsExpected int
	match        bool
	expectErr    bool
}

var _ = Describe("pattern", func() {
	It("Should match patterns properly", func() {
		var err error
		Expect(err).To(BeNil())

		tests := []matchTest{
			{
				content:      "set name to wes",
				pattern:      NewPattern("set <string> to <string>", "set ([^ \t\n\r]+) to ([^ \t\n\r]+)", nil),
				argsExpected: 2,
				match:        true,
				expectErr:    false,
			},
			{
				content:      "this should not match",
				pattern:      NewPattern("set <string> to <string>", "set ([^ \t\n\r]+) to ([^ \t\n\r]+)", nil),
				argsExpected: 2,
				match:        false,
				expectErr:    false,
			},
			{
				content:      "increase the volume",
				pattern:      NewPattern("increase the volume", "increase the volume", nil),
				argsExpected: 0,
				match:        true,
				expectErr:    false,
			},
		}

		for _, t := range tests {

			match, args, err := t.pattern.Match(t.content)

			if t.match {
				Expect(len(args)).To(Equal(t.argsExpected))
			}

			if t.expectErr {
				Expect(err).ToNot(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(match).To(Equal(t.match))
			}
		}
	})
})
