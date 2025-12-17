package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func extractRegionAndMsg(str string) (string, string) {
	if strings.Contains(str, ":") {
		colonSplit := strings.Split(str, ":")
		region := strings.TrimSpace(colonSplit[0])
		msg := strings.TrimSpace(colonSplit[1])

		return region, msg
	}

	return "", str
}

func BuildPrefixWithMsg(template Template, prefix string, msg string) string {
	region, msg := extractRegionAndMsg(msg)

	if region != "" {
		return ExpandTemplate(template.Region, prefix, region, msg)
	}

	return ExpandTemplate(template.Normal, prefix, region, msg)
}

func ReplaceHeaderFromCommit(text string, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	defer file.Sync()

	body , err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// split by end of line
	splitByEOL := strings.Split(string(body), "\n")

	// replace the first line with text
	splitByEOL[0] = text

	file.Truncate(0)
	file.WriteAt([]byte(strings.Join(splitByEOL, "\n")), 0)

	return nil
}

func GetCommitMsgFromFile(fileName string) (string, error) {
	// the first line of the file should be the commit msg

	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}

	if scanner := bufio.NewScanner(file); scanner.Scan() {
		line := scanner.Text()
		
		return line, nil
	}


	return "", fmt.Errorf("invalid")

}
