package glossa

import (
	"fmt"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var _ Glossa = (*glossa)(nil)

type Glossa interface {
	Parse(command string) (Pattern, error)
}

type glossa struct {
	tokens   []string
	tokenIds map[string]int
	lexer    *lexmachine.Lexer
}

func (g *glossa) Parse(command string) (Pattern, error) {
	var err error

	s, err := g.lexer.Scanner([]byte(command))

	if err != nil {
		return nil, err
	}

	regexStr := ""
	patternTokens := make([]*lexmachine.Token, 0)

	nonSpaceTokens := 0
	slotTokens := 0

	for tok, err, eof := s.Next(); !eof; tok, err, eof = s.Next() {
		if _, is := err.(*machines.UnconsumedInput); is {
			err = fmt.Errorf("bad token in command")
			return nil, err
		} else if err != nil {
			return nil, err
		}
		token := tok.(*lexmachine.Token)

		patternTokens = append(patternTokens, token)

		switch token.Type {
		case g.tokenIds["ID"]:
			regexStr += string(token.Lexeme)
			nonSpaceTokens++
		case g.tokenIds["INTEGER"]:
			regexStr += `([0-9]+[^\.]{0,1})`
			nonSpaceTokens++
		case g.tokenIds["FLOAT"]:
			regexStr += `([0-9]+\.[0-9]+)`
			nonSpaceTokens++
		case g.tokenIds["SLOT"]:
			regexStr += "([^ \t\n\r]+)"
			nonSpaceTokens++
			slotTokens++
		case g.tokenIds["SPACE"]:
			regexStr += " "
		}
	}

	if nonSpaceTokens == 1 && slotTokens == 1 {
		return nil, fmt.Errorf("invalid pattern: has one token that is not a space and that token is not an ID")
	}

	if nonSpaceTokens == slotTokens {
		return nil, fmt.Errorf("invalid pattern: contains only SLOT tokens")
	}

	pat := NewPattern(command, regexStr, patternTokens)

	return pat, nil
}

func (g *glossa) init() error {
	g.initTokens()
	var err error
	g.lexer, err = g.initLexer()
	return err
}

func (g *glossa) initTokens() {
	g.tokens = []string{
		"ID",
		"INTEGER",
		"FLOAT",
		"SLOT",
		"SPACE",
	}
	g.tokenIds = make(map[string]int)
	for i, tok := range g.tokens {
		g.tokenIds[tok] = i
	}
}

func (g *glossa) initLexer() (*lexmachine.Lexer, error) {
	lexer := lexmachine.NewLexer()

	lexer.Add([]byte(`([a-z]|[A-Z]|_)+`), g.token("ID"))
	lexer.Add([]byte(`[0-9]+\.[0-9]+`), g.token("FLOAT"))
	lexer.Add([]byte(`[0-9]+[^\.]{0,1}`), g.token("INTEGER"))
	lexer.Add([]byte("( |\t|\n|\r)+"), g.token("SPACE"))
	lexer.Add([]byte(`\<`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			str := make([]byte, 0, 10)
			str = append(str, match.Bytes...)
			brackets := 1
			match.EndLine = match.StartLine
			match.EndColumn = match.StartColumn
			for tc := scan.TC; tc < len(scan.Text); tc++ {
				str = append(str, scan.Text[tc])
				match.EndColumn += 1
				if scan.Text[tc] == '\n' {
					match.EndLine += 1
				}
				if scan.Text[tc] == '<' {
					brackets += 1
				} else if scan.Text[tc] == '>' {
					brackets -= 1
				}
				if brackets == 0 {
					match.TC = scan.TC
					scan.TC = tc + 1
					match.Bytes = str[1 : len(str)-1]
					x, _ := g.token("SLOT")(scan, match)
					t := x.(*lexmachine.Token)
					v := t.Value.(string)
					t.Value = v[1 : len(v)-1]
					return t, nil
				}
			}
			return nil,
				fmt.Errorf("unclosed HTML literal starting at %d, (%d, %d)",
					match.TC, match.StartLine, match.StartColumn)
		},
	)

	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}

func (g *glossa) token(name string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(g.tokenIds[name], string(m.Bytes), m), nil
	}
}

func NewParser() (Glossa, error) {
	g := &glossa{}
	err := g.init()
	return g, err
}
