package smhb

import (
	"os"
	"testing"
)

func TestGetPostAddresses(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(TITLE, CONTENT, HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	addresses, err := client.GetPostAddresses(HANDLE, *token)

	if err != nil {
		t.Error(err)
		return
	}

	if len(addresses) != 1 {
		t.Error("post addresses count mismatch")
		return
	}

	_, err = getPost(addresses[0], context, access)

	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetPost(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(TITLE, CONTENT, HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(HANDLE, context)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := client.GetPost(HANDLE, addresses[0], *token)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Title(), TITLE, t)
	match(out.Content(), CONTENT, t)
	match(out.Author(), HANDLE, t)
}

func TestAddPost(t *testing.T) {
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

	err = client.AddPost(TITLE, CONTENT, HANDLE, *token)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(HANDLE, context)

	if err != nil {
		t.Error(err)
		return
	}

	// Get post locally
	out, err := getPost(addresses[0], context, access)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Title(), TITLE, t)
	match(out.Content(), CONTENT, t)
	match(out.Author(), HANDLE, t)
}

func TestEditPost(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(TITLE, CONTENT, HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(HANDLE, context)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditPost(addresses[0], TITLE+"_", CONTENT+"_", *token)

	if err != nil {
		t.Error(err)
		return
	}

	// Get post locally
	out, err := getPost(addresses[0], context, access)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Title(), TITLE+"_", t)
	match(out.Content(), CONTENT+"_", t)
	match(out.Author(), HANDLE, t)
}

func TestDeletePost(t *testing.T) {
	client, context, access := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(HANDLE, PASSWORD, NAME, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(TITLE, CONTENT, HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(HANDLE, context)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.DeletePost(addresses[0], *token)

	if err != nil {
		t.Error(err)
		return
	}

	// Check if post exists

	_, err = getPost(addresses[0], context, access)

	if err == nil {
		t.Error("post found after deletion")
	}
}

func TestGetFeed(t *testing.T) {
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

	// Create test posts

	err = addPost(TITLE, CONTENT, HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	err = addPost(TITLE+"_", CONTENT+"_", HANDLE+"_", context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool locally
	pool, err := getPool(HANDLE, context, access)

	if err != nil {
		t.Error(err)
		return
	}

	// Add user locally
	err = pool.add(HANDLE+"_", context, access)

	if err != nil {
		t.Error(err)
		return
	}

	token, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := client.GetFeed(HANDLE, *token)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = getPost(out.Addresses()[0], context, access)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = getPost(out.Addresses()[1], context, access)

	if err != nil {
		t.Error(err)
		return
	}
}
