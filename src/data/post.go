package data

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
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

	dir := DATA_PATH + "/" + this.Author + "/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	return ioutil.WriteFile(
		path(this.Author+"/"+this.timestamp, "post"), buffer, 0600,
	)
}

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
