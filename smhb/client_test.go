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

func TestAddUser(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	client.AddUser(HANDLE, PASSWORD, NAME)

	// Get user locally
	out, err := getUser(context, HANDLE)

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

func TestGetPool(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	addUser(context, HANDLE, PASSWORD, NAME)

	out, err := client.GetPool(HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	if handle := out.Handle(); handle != HANDLE {
		t.Error(handle, "does not match", HANDLE)
	}
}

func TestEditPoolAdd(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test users
	addUser(context, HANDLE, PASSWORD, NAME)
	addUser(context, HANDLE+"_", PASSWORD, NAME)

	// Get pool locally
	out, err := getPool(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditPoolAdd(HANDLE, HANDLE+"_")

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool and its users locally
	out, err = getPool(context, HANDLE)
	users := out.Users()

	if _, ok := users[HANDLE+"_"]; !ok {
		t.Error("added user not found in pool")
	}
}

func TestGetPostAddresses(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	addUser(context, HANDLE, PASSWORD, NAME)

	// Create test post
	addPost(context, TITLE, CONTENT, HANDLE)

	addresses, err := client.GetPostAddresses(HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	if len(addresses) != 1 {
		t.Error("post addresses count mismatch")
		return
	}

	_, err = getPost(context, addresses[0])

	if err != nil {
		t.Error(err)
		return
	}
}
