package smhb

import (
	"fmt"
	"io/ioutil"
	"log"
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

/* Interface implementation */

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

func (this *post) GetDir() string {
	return this.author + "/"
}

func (this *post) GetPath() string {
	return this.GetAddress() + ".post"
}

func (this *post) String() string {
	return "post \"" + this.title + "\" by \"" + this.author + "\""
}

// Internal post data wrapper for serialization
type postData struct {
	Title     string
	Content   string
	WasEdited bool
}

// Update post title and content
func (this *post) update(
	title, content string, context serverContext, access Access,
) error {
	// Set new data
	this.title = title
	this.content = content

	// Set edited flag
	this.wasEdited = true

	// Save
	err := access.SaveWithDir(this, this.GetDir(), true, context)

	if err != nil {
		return err
	}

	log.Printf("Updated %s", this)

	return nil
}

func (this *post) WasAuthoredBy(handle string) bool {
	return this.author == handle
}

// Get post unique identifier
func (this *post) GetAddress() string {
	return this.GetDir() + this.timestamp
}

// Create new post and save
func addPost(
	title, content, author string, context serverContext, access Access,
) error {
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
	err := access.SaveWithDir(this, this.GetDir(), true, context)

	if err != nil {
		return err
	}

	log.Printf("Created %s", this)

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
func removePost(
	address string, context serverContext, access Access,
) error {
	// Get author and timestamp
	tokens := strings.Split(address, "/")
	author, stamp := tokens[0], tokens[1]

	this := &post{
		author:    author,
		timestamp: stamp,
	}

	if err := access.Remove(this, context); err != nil {
		return err
	}

	return nil
}

// Remove all posts for a user with a given handle
func removePostsByAuthor(
	author string, context serverContext, access Access,
) error {
	if err := access.RemoveDir(author, context); err != nil {
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
	wrapper := &postData{}
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
