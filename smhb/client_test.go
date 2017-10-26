package smhb

import (
	"os"
	"testing"
)

const (
	HANDLE   = "test_handle"
	PASSWORD = "test_password"
	NAME     = "test_name"
	TEST_DIR = "./tmp"
)

func bootstrap() (Client, serverContext) {
	os.Mkdir(TEST_DIR, os.ModePerm)

	server := NewServer("localhost", 19138, TCP, TEST_DIR)
	testContext := serverContext{server.DataPath()}
	go server.ListenAndServe()

	return NewClient("localhost", 19138, TCP), testContext
}

func TestGetUser(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	addUser(context, HANDLE, PASSWORD, NAME)

	out, err := client.GetUser(HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	if handle := out.Handle(); handle != HANDLE {
		t.Error(handle, "does not match", HANDLE)
	}

	if err := out.Validate(PASSWORD); err != nil {
		t.Error(err)
	}

	if name := out.Name(); name != NAME {
		t.Error(name, "does not match", NAME)
	}
}
