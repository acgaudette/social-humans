package data

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// Post data representation structure
type Post struct {
	Title     string
	Content   string
	Author    string
	Timestamp string
}

// Post data wrapper for serialization
type postData struct {
	Title   string
	Content string
}

// The address is a unique string identifier for the post
func (this *Post) GetAddress() string {
	return this.Author + "/" + this.Timestamp
}

// Write post to file
func (this *Post) save() error {
	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	// Create user directory if it doesn't already exist
	dir := prefix(this.Author + "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	// Write to file
	return ioutil.WriteFile(
		prefix(this.GetAddress()+".post"), buffer, 0600,
	)
}

// Create new post and save
func AddPost(title, content string, author *User) error {
	stamp := time.Now().UTC().Format(TIMESTAMP_LAYOUT)

	this := &Post{
		Title:     title,
		Content:   content,
		Author:    author.Handle,
		Timestamp: stamp,
	}

	// Update data
	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Created post \"%s\" by \"%s\"", title, author)

	return nil
}

// Return addresses of all post files for a given user handle
func GetPostAddresses(author string) ([]string, error) {
	// Read posts directory for user
	files, err := ioutil.ReadDir(prefix(author + "/"))

	if err != nil {
		return nil, err
	}

	addresses := []string{}

	// Build addresses slice
	for _, file := range files {
		// Get address from filename
		address := author + "/" + strings.Split(file.Name(), ".")[0]
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Load post data with lookup address
func LoadPost(address string) (*Post, error) {
	// Read post file
	buffer, err := ioutil.ReadFile(prefix(address + ".post"))

	if err != nil {
		return nil, err
	}

	// Get author and timestamp
	tokens := strings.Split(address, "/")
	author, stamp := tokens[0], tokens[1]

	loaded := &Post{
		Author:    author,
		Timestamp: stamp,
	}

	// Deserialize the rest of the data
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded post \"%s.post\"", address)

	return loaded, nil
}

// Update post title and content
func UpdatePost(address, title, content string) error {
	// Confirm that the post already exists
	post, err := LoadPost(address)

	if err != nil {
		return err
	}

	// Create new post structure with updated title and content
	this := &Post{
		Title:     title,
		Content:   content,
		Author:    post.Author,
		Timestamp: post.Timestamp,
	}

	// Update data
	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Updated post \"%s\" by \"%s\"", title, this.Author)

	return nil
}

// Remove post with lookup address
func RemovePost(address string) error {
	if err := os.Remove(prefix(address + ".post")); err != nil {
		return err
	}

	log.Printf("Deleted post \"%s\"", address)

	return nil
}

/* Satisfy binary interfaces */

func (this *Post) MarshalBinary() ([]byte, error) {
	// Create wrapper from post struct
	wrapper := &postData{
		Title:   this.Title,
		Content: this.Content,
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode wrapper with gob
	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *Post) UnmarshalBinary(buffer []byte) error {
	wrapper := postData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode wrapper with gob
	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	// Load wrapper into new user struct
	this.Title = wrapper.Title
	this.Content = wrapper.Content

	return nil
}
