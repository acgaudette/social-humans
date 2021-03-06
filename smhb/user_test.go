package smhb

import (
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

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
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	client.AddUser(HANDLE, PASSWORD, NAME)

	// Get user locally
	out, err := getUser(HANDLE, context, access)

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
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditUserName(HANDLE, NAME+"_", *token)

	if err != nil {
		t.Error(err)
		return
	}

	// Get user locally
	out, err := getUser(HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Name(), NAME+"_", t)
}

func TestEditUserPassword(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditUserPassword(HANDLE, PASSWORD+"_", *token)

	if err != nil {
		t.Error(err)
		return
	}

	// Get user locally
	out, err := getUser(HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	if err = out.validate(PASSWORD + "_"); err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.DeleteUser(HANDLE, *token)

	if err != nil {
		t.Error(err)
		return
	}

	// Check if user exists
	_, err = getUser(HANDLE, context, access)

	if err == nil {
		t.Error("user found after deletion")
	}
}
