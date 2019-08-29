package fritzbox

import (
	"net/url"

	"github.com/stretchr/testify/mock"
)

type ParserMock struct {
	mock.Mock
}

func (p *ParserMock) ParseQuery(query string) (url.Values, error) {
	args := p.Called(query)
	return args.Get(0).(url.Values), args.Error(1)
}
