package glossa_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGlossa(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Glossa Suite")
}
