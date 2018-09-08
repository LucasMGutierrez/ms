package ms2

import (
	"ms/delay"
	"ms/config"
)

type Args int

type Microservice2 int

func (m *Microservice2) Get(args Args, reply *string) error {
	delay.Sleep(config.Ms2Delay, config.Ms2DelayVar)
	*reply = "Hello"

	return nil
}
