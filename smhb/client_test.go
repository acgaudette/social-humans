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

func match(in, out string, t *testing.T) {
	if in != out {
		t.Error(in, "does not match", out)
	}
}

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

	if err := out.Validate(PASSWORD); err != nil {
		t.Error(err)
	}

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

	if err := out.Validate(PASSWORD); err != nil {
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

	err = client.EditUserName(HANDLE, NAME+"_")

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

	err = client.EditUserPassword(HANDLE, PASSWORD+"_")

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

	if err = out.Validate(PASSWORD+"_"); err != nil {
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

	err = client.DeleteUser(HANDLE)

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

func TestGetPool(t *testing.T) {
	client, context := bootstrap()
	defer os.RemoveAll(TEST_DIR)

	// Create test user
	_, err := addUser(context, HANDLE, PASSWORD, NAME)

	if err != nil {
		t.Error(err)
		return
	}

	out, err := client.GetPool(HANDLE)

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

	err = client.EditPoolBlock(HANDLE, HANDLE+"_")

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

	out, err := client.GetPost(addresses[0])

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

	err = client.AddPost(TITLE, CONTENT, HANDLE)

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

	err = client.EditPost(addresses[0], TITLE+"_", CONTENT+"_")

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

	err = client.DeletePost(addresses[0])

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
