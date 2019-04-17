package config

import (
	"fmt"
	"testing"
)

func setup(t *testing.T) {
	// TODO new file
}

func teardown(t *testing.T) {
	// TODO delete file
}

func TestSetupConfig(t *testing.T) {
	setup(t)
	defer teardown(t)
	c := New()
	err := c.Setup("")
	fmt.Println(err)
}
