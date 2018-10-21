package delay

import (
    "time"
    "math/rand"
    "fmt"
)

func Sleep(sleepTime time.Duration, sleepVar time.Duration, init *bool) {
    fSleepTime := float64(sleepTime)
    fSleepVar := float64(sleepVar)

    if !(*init) {
    	*init = true
    	rand.Seed(time.Now().UnixNano())
    	fmt.Println("ola")
    }

	delay := time.Duration(fSleepTime + (fSleepVar * rand.NormFloat64()))
	time.Sleep(delay)
}
