package smhb

import (
	"os"
	"testing"
)

func TestGetPool(t *testing.T) {
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

	out, err := client.GetPool(HANDLE, *tok)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Handle(), HANDLE, t)
}

func TestEditPoolAdd(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test users

	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = addUser(context, HANDLE+"_", PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool locally
	out, err := getPool(context, HANDLE)

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
	out, err = getPool(context, HANDLE)
	users := out.Users()

	if _, ok := users[HANDLE+"_"]; !ok {
		t.Error("added user not found in pool")
	}
}

func TestEditPoolBlock(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test users

	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = addUser(context, HANDLE+"_", PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool locally
	out, err := getPool(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	// Add user locally
	err = out.add(context, HANDLE+"_")

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
	out, err = getPool(context, HANDLE)
	users := out.Users()

	if _, ok := users[HANDLE+"_"]; ok {
		t.Error("blocked user found in pool")
	}
}
