package smhb

import (
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	out, err := client.GetUser(HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Handle(), HANDLE, t)
	match(out.Name(), NAME, t)
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

	match(out.Handle(), HANDLE, t)

	if err := out.validate(PASSWORD); err != nil {
		t.Error(err)
	}

	match(out.Name(), NAME, t)
}

func TestEditUserName(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditUserName(HANDLE, NAME+"_", *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Get user locally
	out, err := getUser(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Name(), NAME+"_", t)
}

func TestEditUserPassword(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditUserPassword(HANDLE, PASSWORD+"_", *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Get user locally
	out, err := getUser(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	if err = out.validate(PASSWORD + "_"); err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.DeleteUser(HANDLE, *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Check if user exists
	_, err = getUser(context, HANDLE)

	if err == nil {
		t.Error("user found after deletion")
	}
}
