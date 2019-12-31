package main

import (
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(1)

	p.Match("set name to wes", func(parameters ...interface{}) {
		defer wg.Done()

		logrus.WithFields(logrus.Fields{
			"properties_len": len(parameters),
		}).Infof("match result parameters=%+#v", parameters)
	})

	wg.Wait()
}
