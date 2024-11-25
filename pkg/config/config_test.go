package config

import (
	"testing"
)

func TestMustLoadByPath(t *testing.T) {
	cfg := MustLoadByPath("../../config/local.yaml")

	t.Logf("Конфигурация local: %+v", cfg)

	cfg = MustLoadByPath("../../config/local.yaml")

	t.Logf("Конфигурация prod: %+v", cfg)
}
