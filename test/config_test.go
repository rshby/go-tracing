package test

import (
	"fmt"
	"go-tracing/internal/config"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fmt.Println(config.AppPort())

	fmt.Println(config.DatabaseHost())

	fmt.Println(config.DatabaseUser())

	fmt.Println(config.DatabasePassword())

	fmt.Println(config.DatabasePort())

	fmt.Println(config.DatabaseName())

	fmt.Println(config.DatabaseDSN())
}
