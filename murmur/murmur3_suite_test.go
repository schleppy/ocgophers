package murmur3_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMurmur3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Murmur3 Suite")
}
