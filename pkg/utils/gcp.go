package utils

import (
	"fmt"
	"strings"
)

var website = "https://storage.googleapis.com"

func MakeGCPUrl(bucket string, dest string) (url string) {
	url = fmt.Sprintf("%s/%s/%s", website, bucket, dest)
	url = strings.ReplaceAll(url, " ", "%20")
	return url
}

func GetUrlAndFilename(bucket, dest string) (url string, filename string) {
	url = website + "/" + bucket + "/" + dest
	url = strings.ReplaceAll(url, " ", "%20")

	splits := strings.Split(dest, "/")
	filename = splits[len(splits)-1]

	return url, filename
}
