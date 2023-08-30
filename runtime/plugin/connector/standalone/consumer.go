package standalone

import "go.uber.org/atomic"

type consumer struct {
	started atomic.Bool
}
