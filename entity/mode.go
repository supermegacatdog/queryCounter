package entity

import "github.com/supermegacatdog/queryCounter/pkg"

type Mode int

const CounterModeScan = 1
const CounterModeRead = 2

func (m Mode) Validate() error {
	switch m {
	case
		CounterModeScan,
		CounterModeRead:
		return nil

	default:
		return pkg.InvalidMode
	}
}
