package strings

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const STR_LEN = 100

func TestRandStringLength(t *testing.T) {
	var tests = []struct {
		name   string
		strLen int
		expLen int
	}{
		{
			name:   "0 length",
			strLen: 0,
			expLen: 0,
		}, {
			name:   "1 length",
			strLen: 1,
			expLen: 1,
		}, {
			name:   "-1 length",
			strLen: -1,
			expLen: 0,
		}, {
			name:   "100 length",
			strLen: 100,
			expLen: 100,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			assert.Len(t, RandAlphabeticString(testItem.strLen), testItem.expLen)
			assert.Len(t, RandHexAlphaNumericString(testItem.strLen), testItem.expLen)
		})
	}
}

func TestRandAlphabeticalStringRandomness(t *testing.T) {
	max := int(math.Pow(2, 10))
	seen := make([]string, 0)

	for i := 0; i < max; i++ {
		str := RandAlphabeticString(STR_LEN)
		assert.NotSubset(t, seen, []string{str})

		seen = append(seen, str)
	}
}

func TestRandHexAlphaNumericStringRandomness(t *testing.T) {
	max := int(math.Pow(2, 10))
	seen := make([]string, 0)

	for i := 0; i < max; i++ {
		str := RandHexAlphaNumericString(STR_LEN)
		assert.NotSubset(t, seen, []string{str})

		seen = append(seen, str)
	}
}

func TestAppend(t *testing.T) {
	r := RandAlphabeticString(STR_LEN)

	var tests = []struct {
		name string
		og   string
		add  []string
		exp  string
	}{
		{
			name: "empty original",
			og:   "",
			add:  []string{r},
			exp:  r,
		}, {
			name: "original with single",
			og:   "og",
			add:  []string{r},
			exp:  "og" + r,
		}, {
			name: "original with double",
			og:   "og",
			add:  []string{r, r},
			exp:  "og" + r + r,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			assert.Equal(t, Append(testItem.og, testItem.add...), testItem.exp)
		})
	}
}
