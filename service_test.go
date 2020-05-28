package golfs

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestServiceDefaultNoNamespace(t *testing.T) {
	_, err := setupService()
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err != ErrNoNamespace {
		t.Errorf("expected ErrNoNamespace, got: %v", err)
	}
}

func TestServiceDefaultNoGithubToken(t *testing.T) {
	tns := "test"
	os.Setenv("GOLFS_DS_NAMESPACE", tns)
	defer func() {
		os.Unsetenv("GOLFS_DS_NAMESPACE")
	}()

	_, err := setupService()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != ErrNoGithubToken {
		t.Errorf("wrong error, expected ErrNoGithubToken, got '%v'", err)
	}
}

func TestServiceDefaultWithNamespace(t *testing.T) {
	tns := "test"
	os.Setenv("GOLFS_DS_NAMESPACE", tns)
	os.Setenv("GOLFS_GITHUB_TOKEN", "token")
	defer func() {
		os.Unsetenv("GOLFS_DS_NAMESPACE")
		os.Unsetenv("GOLFS_GITHUB_TOKEN")
	}()

	c, err := setupService()
	if err != nil {
		t.Fatalf("expected service, got error: %v", err)
	}

	if c.namespace != tns {
		t.Errorf("wrong namespace, expected '%v' got '%v'", tns, c.namespace)
	}

	shouldTimeout := time.Minute * 5
	shouldLock := true

	if c.lockTimeout != shouldTimeout {
		t.Errorf("expected default timeout of '%v', got '%v'", shouldTimeout.String(), c.lockTimeout.String())
	}

	if c.locking != shouldLock {
		t.Errorf("expected locking enabled value of '%v', got '%v'", shouldLock, c.locking)
	}
}

func TestServiceSetTimeout(t *testing.T) {
	tests := []struct {
		in string
		ex time.Duration
		s  bool
	}{
		{"1m", time.Minute, true},
		{"1h", time.Hour, true},
		{"30s", time.Second * 30, true},
		{"1d", time.Second, false},
		{"nope", time.Second, false},
	}

	tns := "test"
	os.Setenv("GOLFS_DS_NAMESPACE", tns)
	os.Setenv("GOLFS_GITHUB_TOKEN", "token")
	defer func() {
		os.Unsetenv("GOLFS_DS_NAMESPACE")
		os.Unsetenv("GOLFS_GITHUB_TOKEN")
	}()

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test_%v", i), func(t *testing.T) {
			os.Setenv("GOLFS_LOCK_TIMEOUT", tt.in)
			defer func() {
				os.Unsetenv("GOLFS_LOCK_TIMEOUT")
			}()

			c, err := setupService()
			if tt.s {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if c.namespace != tns {
					t.Errorf("wrong namespace, expected '%v', got '%v'", tns, c.namespace)
				}

				if c.lockTimeout != tt.ex {
					t.Errorf("wrong timeout, expected '%v', got '%v'", tt.ex.String(), c.lockTimeout.String())
				}
			} else {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			}
		})
	}
}
