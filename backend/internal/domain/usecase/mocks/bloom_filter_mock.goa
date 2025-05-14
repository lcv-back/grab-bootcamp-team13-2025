package mocks

import (
	"github.com/stretchr/testify/mock"
)

type BloomFilter struct {
	mock.Mock
}

func (m *BloomFilter) Add(data []byte) {
	m.Called(data)
}

func (m *BloomFilter) Test(data []byte) bool {
	args := m.Called(data)
	return args.Bool(0)
}

func (m *BloomFilter) TestAndAdd(data []byte) bool {
	args := m.Called(data)
	return args.Bool(0)
}

func (m *BloomFilter) Reset() {
	m.Called()
}