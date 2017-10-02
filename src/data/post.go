package data

import (
	"io/ioutil"
	"time"
)

type post struct {
	title    string
	content  string
	author   string
	metadata []string
}

func (this *post) save() error {
	name := time.Now().UTC().Format(TIME_LAYOUT)+"_"+this.author
	return ioutil.WriteFile(
		path(name, "post"),
		[]byte(this.content),
		0600,
	)
}

func NewPost(title, content, author string) error {
	this := &post{
		title: title,
		content: content,
		author: author,
	}

	return this.save()
}
