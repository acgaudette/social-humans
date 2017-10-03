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

type Post struct {
	Title     string
	Content   string
	Author    string
	timestamp string
}

type postData struct {
	Title   string
	Content string
}

func (this *Post) MarshalBinary() ([]byte, error) {
	wrapper := &postData{
		Title:   this.Title,
		Content: this.Content,
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *Post) UnmarshalBinary(buffer []byte) error {
	wrapper := postData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	this.Title = wrapper.Title
	this.Content = wrapper.Content

	return nil
}

// Write post to file
func (this *Post) save() error {
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	dir := prefix(this.Author + "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	return ioutil.WriteFile(
		prefix(this.Author+"/"+this.timestamp+".post"), buffer, 0600,
	)
}

// Create new post and save
func NewPost(title, content, author string) error {
	stamp := time.Now().UTC().Format(TIME_LAYOUT)

	this := &Post{
		Title:     title,
		Content:   content,
		Author:    author,
		timestamp: stamp,
	}

	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Created post \"%s\" by \"%s\"", title, author)

	return nil
}

// Return titles of all post files for a given user
func GetPostAddresses(author string) ([]string, error) {
	// Read posts directory for user
	files, err := ioutil.ReadDir(prefix(author + "/"))

	if err != nil {
		return nil, err
	}

	addresses := []string{}

	for _, file := range files {
		// Get address from filename
		address := author + "/" + strings.Split(file.Name(), ".")[0]
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Where the address is "author/timestamp"
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
		timestamp: stamp,
	}

	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded post \"%s.post\"", address)

	return loaded, nil
}

func RemovePost(address string) error {
	if err := os.Remove(prefix(address + ".post")); err != nil {
		return err
	}

	log.Printf("Deleted post \"%s\"", address)

	return nil
}

