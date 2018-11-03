package config

import (
    "time"
)

var (
	// Frontend delay
    FrontendDelay = 100 * time.Millisecond

    // Frontend delay variation
    FrontendDelayVar = FrontendDelay / 40

    // Ms1 delay
    Ms1Delay = 300 * time.Millisecond

    // Ms1 delay variation
    Ms1DelayVar = Ms1Delay / 40

    // Ms2 delay
    Ms2Delay = 1000 * time.Millisecond

    // Ms2 delay variation
	Ms2DelayVar = Ms2Delay / 50

	// Ms3 delay
    Ms3Delay = 500 * time.Millisecond

    // Ms3 delay variation
	Ms3DelayVar = Ms3Delay / 30

	// Port frontend
	PortFrontend = ":1234"

	// Port ms1
	PortMs1 = ":1235"

	// Port ms2
	PortMs2 = ":1236"

	// Port ms3
	PortMs3 = ":1237"
)
