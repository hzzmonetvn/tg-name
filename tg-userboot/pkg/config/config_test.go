package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY_PLAINUB", "custom_val")
	defer os.Unsetenv("TEST_KEY_PLAINUB")

	val := getEnv("TEST_KEY_PLAINUB", "default_val")
	if val != "custom_val" {
		t.Errorf("expected 'custom_val', got %s", val)
	}

	valDef := getEnv("NON_EXISTENT_KEY", "default_val")
	if valDef != "default_val" {
		t.Errorf("expected 'default_val', got %s", valDef)
	}
}

func TestConfigLoad(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.CmdTrigger == "" {
		t.Errorf("expected CmdTrigger to have a default value")
	}
}
