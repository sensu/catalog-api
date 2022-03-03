// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	catalogv1 "github.com/sensu/catalog-api/internal/api/catalog/v1"
	integrationloader "github.com/sensu/catalog-api/internal/integrationloader"

	mock "github.com/stretchr/testify/mock"
)

// Loader is an autogenerated mock type for the Loader type
type Loader struct {
	mock.Mock
}

// GetFileContentsAsBytes provides a mock function with given fields: _a0
func (_m *Loader) GetFileContentsAsBytes(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFileContentsAsString provides a mock function with given fields: _a0
func (_m *Loader) GetFileContentsAsString(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadChangelog provides a mock function with given fields:
func (_m *Loader) LoadChangelog() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadConfig provides a mock function with given fields:
func (_m *Loader) LoadConfig() (catalogv1.Integration, error) {
	ret := _m.Called()

	var r0 catalogv1.Integration
	if rf, ok := ret.Get(0).(func() catalogv1.Integration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(catalogv1.Integration)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadImages provides a mock function with given fields:
func (_m *Loader) LoadImages() (integrationloader.Images, error) {
	ret := _m.Called()

	var r0 integrationloader.Images
	if rf, ok := ret.Get(0).(func() integrationloader.Images); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(integrationloader.Images)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadLogo provides a mock function with given fields:
func (_m *Loader) LoadLogo() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadReadme provides a mock function with given fields:
func (_m *Loader) LoadReadme() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadResources provides a mock function with given fields:
func (_m *Loader) LoadResources() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}