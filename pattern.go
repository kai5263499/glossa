package glossa

import (
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/timtadh/lexmachine"
)

var _ Pattern = (*pattern)(nil)

type Pattern interface {
	Match(string, func(...interface{})) (bool, error)
	GetCommand() string
	GetTokens() []*lexmachine.Token
	GetRegexStr() string
	GetRe() *regexp.Regexp
}

type pattern struct {
	command  string
	regexStr string
	tokens   []*lexmachine.Token
	re       *regexp.Regexp
}

func (p *pattern) Match(content string, callback func(...interface{})) (bool, error) {
	matches := p.re.FindStringSubmatch(content)
	matchFound := len(matches) > 0
	args := make([]interface{}, 0)

	if matchFound {
		for _, match := range matches[1:] {
			args = append(args, match)
		}
	}

	logrus.WithFields(logrus.Fields{
		"content":     content,
		"regexStr":    p.regexStr,
		"args_len":    len(args),
		"matches len": len(matches),
		"matchFound":  matchFound,
	}).Debugf("match=%+#v args=%+#v", matches, args)

	if matchFound {
		go callback(args...)
	}

	return matchFound, nil
}

func (p *pattern) GetCommand() string {
	return p.command
}

func (p *pattern) GetTokens() []*lexmachine.Token {
	return p.tokens
}

func (p *pattern) GetRegexStr() string {
	return p.regexStr
}

func (p *pattern) GetRe() *regexp.Regexp {
	return p.re
}

func NewPattern(command string, regexStr string, tokens []*lexmachine.Token) Pattern {
	pat := &pattern{
		command:  command,
		regexStr: regexStr,
		re:       regexp.MustCompile(regexStr),
		tokens:   tokens,
	}
	return pat
}
