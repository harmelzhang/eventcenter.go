package standalone

import (
	"go.uber.org/atomic"
)

type producer struct {
	started atomic.Bool
}
