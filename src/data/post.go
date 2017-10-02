package data

import (
	"io/ioutil"
	"log"
	"time"
)

type post struct {
	title     string
	content   string
	author    string
	timestamp string
}

func (this *post) save() error {
	return ioutil.WriteFile(
		path(this.timestamp+"_"+this.author, "post"),
		[]byte(this.content),
		0600,
	)
}

func NewPost(title, content, author string) error {
	stamp := time.Now().UTC().Format(TIME_LAYOUT)

	this := &post{
		title:     title,
		content:   content,
		author:    author,
		timestamp: stamp,
	}

	log.Printf("Created post \"%s\" by \"%s\"", title, author)

	return this.save()
}
