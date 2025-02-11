package config

import (
	"testing"
)

func TestRead(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Errorf("Some stupid error %v", err)
	}
	if config.DbUrl == "" {
		t.Errorf("DbUrl is blank")
	}
}

func TestSetUser(t *testing.T) {
	user := "MatthewNorth"
	config, _ := Read()
	config.SetUser(user)
	if config.CurrentUserName != user {
		t.Errorf("User was not updated to %s, instead it is %s", user, config.CurrentUserName)
	}
	secondRead, _ := Read()
	if secondRead.CurrentUserName != user {
		t.Errorf("Reading the file again, it was not updated to %s, instead it is %s", user, secondRead.CurrentUserName)
	}
}
