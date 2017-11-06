package smhb

import (
	"os"
	"testing"
)

func TestGetPostAddresses(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(context, TITLE, CONTENT, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	addresses, err := client.GetPostAddresses(HANDLE, *tok)

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

func TestGetPost(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(context, TITLE, CONTENT, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := client.GetPost(HANDLE, addresses[0], *tok)

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Title(), TITLE, t)
	match(out.Content(), CONTENT, t)
	match(out.Author(), HANDLE, t)
}

func TestAddPost(t *testing.T) {
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

	err = client.AddPost(TITLE, CONTENT, HANDLE, *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	// Get post locally
	out, err := getPost(context, addresses[0])

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Title(), TITLE, t)
	match(out.Content(), CONTENT, t)
	match(out.Author(), HANDLE, t)
}

func TestEditPost(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(context, TITLE, CONTENT, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.EditPost(addresses[0], TITLE+"_", CONTENT+"_", *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Get post locally
	out, err := getPost(context, addresses[0])

	if err != nil {
		t.Error(err)
		return
	}

	match(out.Title(), TITLE+"_", t)
	match(out.Content(), CONTENT+"_", t)
	match(out.Author(), HANDLE, t)
}

func TestDeletePost(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	// Create test post
	err = addPost(context, TITLE, CONTENT, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	// Get addresses locally
	addresses, err := getPostAddresses(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.DeletePost(addresses[0], *tok)

	if err != nil {
		t.Error(err)
		return
	}

	// Check if post exists

	_, err = getPost(context, addresses[0])

	if err == nil {
		t.Error("post found after deletion")
	}
}

func TestGetFeed(t *testing.T) {
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

	// Create test posts

	err = addPost(context, TITLE, CONTENT, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	err = addPost(context, TITLE+"_", CONTENT+"_", HANDLE+"_")

	if err != nil {
		t.Error(err)
		return
	}

	// Get pool locally
	pool, err := getPool(context, HANDLE)

	if err != nil {
		t.Error(err)
		return
	}

	// Add user locally
	err = pool.add(context, HANDLE+"_")

	if err != nil {
		t.Error(err)
		return
	}

	tok, err := getBackendToken(client, HANDLE, PASSWORD)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := client.GetFeed(HANDLE, *tok)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = getPost(context, out.Addresses()[0])

	if err != nil {
		t.Error(err)
		return
	}

	_, err = getPost(context, out.Addresses()[1])

	if err != nil {
		t.Error(err)
		return
	}
}
