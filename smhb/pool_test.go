package smhb

import (
	"os"
	"testing"
)

func TestGetPool(t *testing.T) {
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

	out, err := client.GetPool(HANDLE, *token)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Handle(), HANDLE, t)
	match(out.Users()[HANDLE], HANDLE, t)
}

func TestEditPoolAdd(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test users

	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = addUser(HANDLE+"_", PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool locally
	out, err := getPool(HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditPoolAdd(HANDLE, HANDLE+"_", *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool and its users locally
	out, err = getPool(HANDLE, context, access)
	users := out.Users()

	if _, ok := users[HANDLE+"_"]; !ok {
		t.Error("added user not found in pool")
	}
}

func TestEditPoolBlock(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test users

	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = addUser(HANDLE+"_", PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool locally
	out, err := getPool(HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Add user locally
	err = out.add(HANDLE+"_", context, access)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditPoolBlock(HANDLE, HANDLE+"_", *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool and its users locally
	out, err = getPool(HANDLE, context, access)
	users := out.Users()

	if _, ok := users[HANDLE+"_"]; ok {
		t.Error("blocked user found in pool")
	}
}
