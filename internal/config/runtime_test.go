package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"SESSION_KEY_FILE", "SESSION_SECRET", "COOKIE_SECURE", "LISTEN_HOST", "PORT", "ENVIRONMENT"} {
		t.Setenv(k, "")
	}
}
func TestRuntimeRejectsAbsentOrPlaceholderKey(t *testing.T) {
	for _, key := range []string{"", "change-this-secret-in-production-minimum-32-chars", strings.Repeat("00", 32)} {
		t.Run(key, func(t *testing.T) {
			resetEnv(t)
			t.Setenv("SESSION_SECRET", key)
			if _, err := Load(); err == nil {
				t.Fatal("insecure key accepted")
			}
		})
	}
}
func TestRuntimePrivateKeyAndProductionDefaults(t *testing.T) {
	resetEnv(t)
	p := filepath.Join(t.TempDir(), "key")
	key := []byte("0123456789abcdef0123456789abcdef")
	if err := os.WriteFile(p, key, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SESSION_KEY_FILE", p)
	t.Setenv("ENVIRONMENT", "production")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !c.SecureCookies || c.ListenAddress != "127.0.0.1:28090" {
		t.Fatalf("bad defaults: %+v", c)
	}
	t.Setenv("COOKIE_SECURE", "false")
	if _, err = Load(); err == nil {
		t.Fatal("insecure production cookies accepted")
	}
}
