package ms3

import (
	"ms/delay"
	"ms/config"
)

type Args int

type Microservice3 int

func (m *Microservice3) Get(args Args, reply *string) error {
	// TODO ms1 + ms2
	delay.Sleep(config.Ms3Delay, config.Ms3DelayVar)
	*reply = "World"

	return nil
}
