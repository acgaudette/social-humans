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

func SavePost(content string, title string, user string) error {
	return ioutil.WriteFile(
		postPath(time.Now().UTC().Format(TIME_LAYOUT))+"_"+title+"_"+user+".post",
		[]byte(content),
		0600,
	)
}

func postPath(handle string) string {
	return DATA_PATH + "/" + handle
}
