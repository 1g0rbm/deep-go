package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	di map[string]func() interface{}
}

func NewContainer() *Container {
	return &Container{
		di: make(map[string]func() interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	fn, ok := constructor.(func() interface{})
	if ok {
		c.di[name] = fn
	}
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	fn, ok := constructor.(func() interface{})
	if !ok {
		return
	}

	var instance interface{}
	c.di[name] = func() interface{} {
		if instance == nil {
			instance = fn()
		}

		return instance
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	dep, ok := c.di[name]
	if !ok {
		return nil, errors.New("no deps")
	}

	return dep(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)

	assert.False(t, u1 == u2)

	container.RegisterSingletonType("MessageService", func() interface{} {
		return &MessageService{}
	})

	singleMessageService1, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	singleMessageService2, err := container.Resolve("MessageService")
	assert.NoError(t, err)

	singleM1 := singleMessageService1.(*MessageService)
	singleM2 := singleMessageService2.(*MessageService)

	assert.True(t, singleM1 == singleM2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
