package test

import (
	"go-web-demo/config"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config.Reader.ReadConfig()

}
