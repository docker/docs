package bootstrap

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBootstrap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bootstrap Suite")
}

var _ = Describe("SplitOnAnd", func() {
	It("should work", func() {
		for _, testCase := range []struct {
			str      string
			expected []string
		}{
			{
				str:      "a",
				expected: []string{"a"},
			},
			{
				str:      "a&b",
				expected: []string{"a", "b"},
			},
			{
				str:      "a&b&c",
				expected: []string{"a", "b", "c"},
			},
			{
				str:      "&&",
				expected: []string{"", "", ""},
			},
			{
				str:      `\&`,
				expected: []string{"&"},
			},
			{
				str:      `\&\&`,
				expected: []string{"&&"},
			},
			{
				str:      `a\&&b`,
				expected: []string{"a&", "b"},
			},
			{
				str:      `a\&ab&`,
				expected: []string{"a&ab", ""},
			},
			{
				str:      `&\&`,
				expected: []string{"", "&"},
			},
		} {
			actual := splitOnAnd(testCase.str)
			Expect(actual).To(Equal(testCase.expected), "test case: %s", testCase.str)
		}
	})
})
