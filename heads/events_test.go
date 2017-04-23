package heads

import (
	"github.com/stretchr/testify/mock"
)

type testEvents struct {
	mock.Mock
}

func (testEvents testEvents) update(event APIEvents) {
	testEvents.Called()
}
