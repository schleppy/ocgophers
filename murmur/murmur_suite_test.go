package murmur_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMurmur(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Murmur Suite")
}
