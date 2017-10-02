package data

import (
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

func (this *Post) save() error {
	return ioutil.WriteFile(
		path(this.timestamp+"_"+this.Author, "post"),
		[]byte(this.Content),
		0600,
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

	log.Printf("Created post \"%s\" by \"%s\"", title, author)

	return this.save()
}
