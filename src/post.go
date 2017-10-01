package main

import (
	"io/ioutil"
	"time"
)

type post struct {
	content  string
	title    string
	metadata []string
}

func savePost(content string, title string, user string) error {
	return ioutil.WriteFile(
		postpath(time.Now().UTC().Format(TIME_LAYOUT))+"_"+title+"_"+user+".post",
		[]byte(content),
		0600,
	)
}

func postpath(handle string) string {
	return DATA_PATH + "/" + handle
}
