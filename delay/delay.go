package delay

import (
    "time"
    "math/rand"
)

func Sleep(sleepTime time.Duration, sleepVar time.Duration) {
    fSleepTime := float64(sleepTime)
    fSleepVar := float64(sleepVar)

	delay := time.Duration(fSleepTime + (fSleepVar * rand.NormFloat64() * 0.1))
	time.Sleep(delay)
}
