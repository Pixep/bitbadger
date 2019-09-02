package bitbadger

import (
	"testing"
)

func TestSetConfig(t *testing.T) {
	username := "Robert"
	password := "azerty"

	SetConfig(Config{
		Username: username,
		Password: password,
	})

	storedConfig := GetConfig()
	if storedConfig.Username != username || storedConfig.Password != password {
		t.Errorf("SetConfig: Username or password mismatch")
	}
}
