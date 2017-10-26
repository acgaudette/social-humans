package smhb

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
type post struct {
	title     string
	content   string
	author    string
	timestamp string
	wasEdited bool
}

/* Interface implementation getters */

func (this *post) Title() string {
	return this.title
}

func (this *post) Content() string {
	return this.content
}

func (this *post) Author() string {
	return this.author
}

func (this *post) Timestamp() string {
	return this.timestamp
}

func (this *post) WasEdited() bool {
	return this.wasEdited
}

// Internal post data wrapper for serialization
type postData struct {
	Title     string
	Content   string
	WasEdited bool
}

// Update post title and content
func (this *post) update(title, content string) error {
	// Set new data
	this.title = title
	this.content = content

	// Set edited flag
	this.wasEdited = true

	// Save
	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Updated post \"%s\" by \"%s\"", title, this.author)

	return nil
}

func (this *post) WasAuthoredBy(handle string) bool {
	return this.author == handle
}

// Get post unique identifier
func (this *post) GetAddress() string {
	return BuildPostAddress(this.author, this.timestamp)
}

// Write post to file
func (this *post) save() error {
	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	// Create user directory if it doesn't already exist
	dir := prefix(this.author + "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	// Write to file
	return ioutil.WriteFile(
		prefix(this.GetAddress()+".post"), buffer, 0600,
	)
}

// Create new post and save
func addPost(title, content, author string) error {
	stamp := time.Now().UTC().Format(TIMESTAMP_LAYOUT)

	this := &post{
		title:     title,
		content:   content,
		author:    author,
		timestamp: stamp,
		wasEdited: false,
	}

	// Update data
	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Created post \"%s\" by \"%s\"", title, author)

	return nil
}

// Load post raw buffer with lookup handle
func loadPost(address string) ([]byte, error) {
	// Read post file
	buffer, err := ioutil.ReadFile(prefix(address + ".post"))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded post \"%s.post\"", address)

	return buffer, nil
}

// Deserialize raw buffer with lookup handle
func deserializePost(address string, buffer []byte) (*post, error) {
	// Get author and timestamp
	tokens := strings.Split(address, "/")
	author, stamp := tokens[0], tokens[1]

	loaded := &post{
		author:    author,
		timestamp: stamp,
	}

	// Deserialize the rest of the data
	err := loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

// Load post data with lookup address
func getPost(address string) (*post, error) {
	buffer, err := loadPost(address)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePost(address, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

// Remove post with lookup address
func removePost(address string) error {
	if err := os.Remove(prefix(address + ".post")); err != nil {
		return err
	}

	log.Printf("Deleted post \"%s\"", address)

	return nil
}

// Remove all posts for a user with a given handle
func removePostsByAuthor(author string) error {
	if err := os.RemoveAll(prefix(author)); err != nil {
		return err
	}

	log.Printf("Deleted all posts by user \"%s\"", author)

	return nil
}

/* Satisfy binary interfaces */

func (this *post) MarshalBinary() ([]byte, error) {
	// Create wrapper from post struct
	wrapper := &postData{
		Title:     this.title,
		Content:   this.content,
		WasEdited: this.wasEdited,
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode wrapper with gob
	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *post) UnmarshalBinary(buffer []byte) error {
	wrapper := postData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode wrapper with gob
	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	// Load wrapper into new user struct
	this.title = wrapper.Title
	this.content = wrapper.Content
	this.wasEdited = wrapper.WasEdited

	return nil
}
