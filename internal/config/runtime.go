package config

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
)

// Runtime contains startup security settings. No signing key has a default.
type Runtime struct {
	SessionKey    []byte
	SecureCookies bool
	ListenAddress string
}

func Load() (Runtime, error) {
	c := Runtime{SecureCookies: os.Getenv("ENVIRONMENT") == "production"}
	path, encoded := os.Getenv("SESSION_KEY_FILE"), os.Getenv("SESSION_SECRET")
	if path != "" && encoded != "" {
		return c, fmt.Errorf("configure only one of SESSION_KEY_FILE and SESSION_SECRET")
	}
	var err error
	if path != "" {
		c.SessionKey, err = os.ReadFile(path)
	} else {
		c.SessionKey, err = hex.DecodeString(encoded)
	}
	if err != nil || len(c.SessionKey) != 32 {
		return c, fmt.Errorf("session signing key must contain exactly 32 random bytes (SESSION_KEY_FILE), or 64 hex characters (SESSION_SECRET)")
	}
	// Reject trivial placeholders without disclosing the supplied key.
	identical := true
	for _, b := range c.SessionKey {
		identical = identical && b == c.SessionKey[0]
	}
	if identical {
		return c, fmt.Errorf("session signing key must not be a repeated-byte placeholder")
	}
	if value := os.Getenv("COOKIE_SECURE"); value != "" {
		c.SecureCookies, err = strconv.ParseBool(value)
		if err != nil {
			return c, fmt.Errorf("COOKIE_SECURE must be a boolean")
		}
	}
	if os.Getenv("ENVIRONMENT") == "production" && !c.SecureCookies {
		return c, fmt.Errorf("production requires secure cookies")
	}
	host, port := os.Getenv("LISTEN_HOST"), os.Getenv("PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "28090"
	}
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return c, fmt.Errorf("PORT must be between 1 and 65535")
	}
	if net.ParseIP(host) == nil {
		return c, fmt.Errorf("LISTEN_HOST must be an IP address")
	}
	c.ListenAddress = net.JoinHostPort(host, port)
	return c, nil
}
