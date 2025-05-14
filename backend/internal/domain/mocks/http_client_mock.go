package mocks

import (
	"grab-bootcamp-be-team13-2025/pkg/utils/http"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(path string, body interface{}) (*http.Response, error) {
	args := m.Called(path, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
} 