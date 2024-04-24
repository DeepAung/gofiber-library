package utils

import (
	"fmt"
	"strings"
)

var website = "https://storage.googleapis.com"

func MakeGCPUrl(bucket string, dest string) string {
	return fmt.Sprintf("%s/%s/%s", website, bucket, dest)
}

func GetUrlAndFilename(bucket, dest string) (url string, filename string) {
	url = website + "/" + bucket + "/" + dest
	splits := strings.Split(dest, "/")
	filename = splits[len(splits)-1]

	return url, filename
}
