package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) GetProductOptionValue(id uint64) (*models.ProductOptionValue, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ProductOptionValue), args.Error(1)
}

func (m *MockDB) CreateProductOptionValue(nu *models.ProductOptionValue) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateProductOptionValue(updated *models.ProductOptionValue) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteProductOptionValue(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0).(time.Time), args.Error(1)
}
