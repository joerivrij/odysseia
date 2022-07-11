package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDurationDifference(t *testing.T) {
	valid := 10
	daysOfFuturePast := valid - 2*valid + 1

	t.Run("RedeployNeeded", func(t *testing.T) {
		created := time.Now().UTC().AddDate(0, 0, daysOfFuturePast).Format(timeFormat)
		redeploy, err := calculateTimeDifference(valid, created)
		assert.Nil(t, err)
		assert.True(t, redeploy)
	})
	t.Run("NoRedeployNeeded", func(t *testing.T) {
		created := time.Now().UTC().Format(timeFormat)
		redeploy, err := calculateTimeDifference(valid, created)
		assert.Nil(t, err)
		assert.False(t, redeploy)
	})
}
