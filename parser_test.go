package glossa

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type parseTest struct {
	command   string
	pattern   Pattern
	expectErr bool
}

var _ = Describe("parser", func() {
	It("Should parse a basic command into a pattern", func() {
		var err error
		Expect(err).To(BeNil())

		g, err := NewParser()
		Expect(err).To(BeNil())
		Expect(g).To(Not(BeNil()))

		pat, err := g.Parse("set <string> to <string>")
		Expect(err).To(BeNil())
		Expect(pat).To(Not(BeNil()))

		matched, args, err := pat.Match("set name to wes")

		Expect(len(args)).To(Equal(2))
		Expect(err).To(BeNil())
		Expect(matched).To(BeTrue())
	})
	It("Should parse commands into corresponding patterns", func() {
		var err error
		Expect(err).To(BeNil())

		g, err := NewParser()
		Expect(err).To(BeNil())
		Expect(g).To(Not(BeNil()))

		tests := []parseTest{
			{command: "play <string>", expectErr: false},
			{command: "play <str", expectErr: true},
			{command: "play <string<>", expectErr: true},
			{command: "play str>", expectErr: true},
			{command: "<string><string>", expectErr: true},
			{command: "<string> <string>", expectErr: true},
			{command: "<string>", expectErr: true},
			{command: "<string> ", expectErr: true},
			{command: " <string> ", expectErr: true},
			{command: "<float>", expectErr: true},
			{command: " <float>", expectErr: true},
			{command: " <float> ", expectErr: true},
			{command: "<float> ", expectErr: true},
			{command: "<integer>", expectErr: true},
			{command: "check <string>", expectErr: false},
			{command: "open <string>", expectErr: false},
			{command: "record", expectErr: false},
			{command: "remind <string> to <string> in <string>", expectErr: false},
		}

		for _, t := range tests {
			p, err := g.Parse(t.command)
			if t.expectErr {
				Expect(err).ToNot(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(p).To(Not(BeNil()))
				Expect(p.GetRe()).To(Not(BeNil()))
				Expect(p.GetTokens()).To(Not(BeNil()))
				Expect(t.command).To(Equal(p.GetCommand()))
			}
		}
	})
})
