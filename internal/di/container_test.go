package di

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer()

	assert.NotNil(t, container, "Expected container to be created")
}

func TestContainer_Add(t *testing.T) {
	container := NewContainer()

	// Test adding first dependency
	err := container.Add("test_key", "test_value")
	assert.NoError(t, err, "Expected no error")

	// Test adding duplicate dependency
	err = container.Add("test_key", "another_value")
	assert.Error(t, err, "Expected error when adding duplicate key")
}

func TestContainer_Get(t *testing.T) {
	container := NewContainer()

	// Test getting non-existent key
	value := container.Get("non_existent")
	assert.Nil(t, value, "Expected nil for non-existent key")

	// Test adding and getting value
	testValue := "test_value"
	err := container.Add("test_key", testValue)
	assert.NoError(t, err, "Expected no error")

	value = container.Get("test_key")
	assert.Equal(t, testValue, value, "Expected %v, got %v", testValue, value)
}

func TestContainer_Concurrency(t *testing.T) {
	container := NewContainer()

	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			key := fmt.Sprintf("key_%d", id)
			value := fmt.Sprintf("value_%d", id)

			err := container.Add(key, value)
			assert.NoError(t, err, "Expected no error")

			retrieved := container.Get(key)
			assert.Equal(t, value, retrieved, "Expected %v, got %v", value, retrieved)

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestContainer_Interface(t *testing.T) {
	var container ContainerInterface
	container = NewContainer()

	assert.NotNil(t, container, "Expected container to implement ContainerInterface")
}
