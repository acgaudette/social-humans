package data

import (
	"io/ioutil"
	"time"
)

type post struct {
	content  string
	title    string
	metadata []string
}

func SavePost(content, title, user string) error {
	return ioutil.WriteFile(
		path(time.Now().UTC().Format(TIME_LAYOUT)+"_"+user, "post"),
		[]byte(content),
		0600,
	)
}
