package oregontrail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomIntRequest(t *testing.T) {
	t.Run("returns a random number between 1 and 10", func(t *testing.T) {
		result := GetRandomInt()

		assert.GreaterOrEqual(t, result, 1, fmt.Sprintf("result %d should be >= 1", result))
		assert.LessOrEqual(t, result, 10, fmt.Sprintf("result %d should be <= 10", result))
	})
}
