package main

import (
	"fmt"
)

var TEST_SERVER_END_POINT string

func init() {
	config := GetConfig()
	TEST_SERVER_END_POINT = fmt.Sprintf(
		"http://%s:%d/",
		config.GetHttpHostname(),
		config.GetHttpPort(),
	)
}
