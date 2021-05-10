package utils

import "strings"

func ExtractFilenameOnly(filepath string) string {
	splits := strings.Split(filepath, "/")
	return splits[len(splits)-1]
}
