package smhb

import (
	"fmt"
	"os"
	"testing"
)

const (
	TEST_DIR = "./tmp"
	HANDLE   = "test_handle"
	PASSWORD = "test_password"
	NAME     = "test_name"
	TITLE    = "test_title"
	CONTENT  = "test_content"
)

func bootstrap() (Client, ServerContext, Access) {
	os.Mkdir(TEST_DIR, os.ModePerm)
	fmt.Fprintf(os.Stderr, "\nBOOTSTRAP\n")

	server := NewServer("localhost", 19138, TCP, 8, TEST_DIR)
	testContext := ServerContext{server.DataPath()}
	go server.ListenAndServe()

	return NewClient("localhost", 19138, TCP), testContext, server.access
}

func getBackendToken(client Client, handle, password string) (*Token, error) {
	token, err := client.GetToken(HANDLE, PASSWORD)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func match(in, out string, t *testing.T) {
	if in != out {
		t.Error("\"" + in + "\" does not match \"" + out + "\"")
	}
}
