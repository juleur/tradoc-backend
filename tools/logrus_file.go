package tools

import (
	"os"

	"github.com/sirupsen/logrus"
)

func CreateLogFile(logPath string) *os.File {
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	return file
}
