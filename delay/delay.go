package delay

import (
    "time"
    "math/rand"
)

func Sleep(sleepTime time.Duration, sleepVar time.Duration) {
    fSleepTime := float64(sleepTime)
    fSleepVar := float64(sleepVar)

	delay := time.Duration((fSleepTime - fSleepVar) + (2 * fSleepVar * rand.Float64()))
	time.Sleep(delay)
}
