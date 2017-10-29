package smhb

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
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
func (this *post) update(context serverContext, title, content string) error {
	// Set new data
	this.title = title
	this.content = content

	// Set edited flag
	this.wasEdited = true

	// Save
	if err := this.save(context); err != nil {
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
	return this.author + "/" + this.timestamp
}

// Write post to file
func (this *post) save(context serverContext) error {
	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	// Create user directory if it doesn't already exist
	dir := prefix(context, this.author+"/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	// Write to file
	return ioutil.WriteFile(
		prefix(context, this.GetAddress()+".post"), buffer, 0600,
	)
}

// Create new post and save
func addPost(context serverContext, title, content, author string) error {
	// Check character limits

	if utf8.RuneCountInString(title) > TITLE_LIMIT {
		return fmt.Errorf("post title length is over %d chars", TITLE_LIMIT)
	}

	if utf8.RuneCountInString(content) > CONTENT_LIMIT {
		return fmt.Errorf("post content length is over %d chars", CONTENT_LIMIT)
	}

	// Create timestamp
	stamp := time.Now().UTC().Format(TIMESTAMP_LAYOUT)

	this := &post{
		title:     title,
		content:   content,
		author:    author,
		timestamp: stamp,
		wasEdited: false,
	}

	// Update data
	if err := this.save(context); err != nil {
		return err
	}

	log.Printf("Created post \"%s\" by \"%s\"", title, author)

	return nil
}

// Load post raw buffer with lookup handle
func loadPost(context serverContext, address string) ([]byte, error) {
	// Read post file
	buffer, err := ioutil.ReadFile(prefix(context, address+".post"))

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
func getPost(context serverContext, address string) (*post, error) {
	buffer, err := loadPost(context, address)

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
func removePost(context serverContext, address string) error {
	if err := os.Remove(prefix(context, address+".post")); err != nil {
		return err
	}

	log.Printf("Deleted post \"%s\"", address)

	return nil
}

// Remove all posts for a user with a given handle
func removePostsByAuthor(context serverContext, author string) error {
	if err := os.RemoveAll(prefix(context, author)); err != nil {
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

	return serialize(wrapper)
}

func (this *post) UnmarshalBinary(buffer []byte) error {
	wrapper := postData{}
	err := deserialize(wrapper, buffer)

	if err != nil {
		return err
	}

	// Load wrapper into new user struct
	this.title = wrapper.Title
	this.content = wrapper.Content
	this.wasEdited = wrapper.WasEdited

	return nil
}
