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

func bootstrap() (Client, serverContext) {
	os.Mkdir(TEST_DIR, os.ModePerm)
	fmt.Fprintf(os.Stderr, "\nBOOTSTRAP\n")

	server := NewServer("localhost", 19138, TCP, 8, TEST_DIR)
	testContext := serverContext{server.DataPath()}
	go server.ListenAndServe()

	return NewClient("localhost", 19138, TCP), testContext
}

func getBackendToken(client Client, handle, password string) (*Token, error) {
	tok, err := client.GetToken(HANDLE, PASSWORD)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func match(in, out string, t *testing.T) {
	if in != out {
		t.Error(in, "does not match", out)
	}
}
