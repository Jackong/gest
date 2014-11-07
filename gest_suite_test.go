package gest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gest Suite")
}
