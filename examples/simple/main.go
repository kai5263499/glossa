package main

import (
	"github.com/kai5263499/glossa"
	"github.com/sirupsen/logrus"
)

func checkErr(msg string, err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Errorf(msg)
		panic(err)
	}
}

func main() {
	var err error

	g, err := glossa.NewParser()
	checkErr("new glossa", err)

	p, err := g.Parse(`set <string> to <string>`)
	checkErr("parse command", err)

	matched, args, err := p.Match("set name to wes")
	checkErr("match", err)

	logrus.WithFields(logrus.Fields{
		"matched":        matched,
		"command":        p.GetCommand(),
		"regexStr":       p.GetRegexStr(),
		"token_len":      len(p.GetTokens()),
		"properties_len": len(args),
	}).Infof("match result parameters=%+#v", args)

}
