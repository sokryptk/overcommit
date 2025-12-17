package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
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

func BuildCommitMessage(template Template, prefix string, scope string, msg string) string {
	if scope != "" {
		return ExpandTemplate(template.Region, prefix, scope, msg)
	}
	return ExpandTemplate(template.Normal, prefix, "", msg)
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
	defer file.Close()

	if scanner := bufio.NewScanner(file); scanner.Scan() {
		line := scanner.Text()
		
		return line, nil
	}


	return "", fmt.Errorf("invalid")
}

func GetScopes() []string {
	seen := make(map[string]bool)
	var scopes []string

	addScope := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			scopes = append(scopes, s)
		}
	}

	// Get scopes from past commits
	out, err := exec.Command("git", "log", "--oneline", "-100", "--format=%s").Output()
	if err == nil {
		re := regexp.MustCompile(`^\w+\(([^)]+)\):`)
		for _, line := range strings.Split(string(out), "\n") {
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				addScope(matches[1])
			}
		}
	}

	// Get top-level directories
	entries, err := os.ReadDir(".")
	if err == nil {
		exclude := map[string]bool{".git": true, "node_modules": true, "vendor": true, ".idea": true, ".vscode": true}
		for _, e := range entries {
			if e.IsDir() && !exclude[e.Name()] && !strings.HasPrefix(e.Name(), ".") {
				addScope(e.Name())
			}
		}
	}

	return scopes
}

func GetStagedDiff() (string, error) {
	out, _ := exec.Command("git", "diff", "--cached", "-p", "--no-color").Output()
	if len(out) > 0 {
		return string(out), nil
	}

	files, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, f := range strings.Split(strings.TrimSpace(string(files)), "\n") {
		if f == "" {
			continue
		}
		content, _ := exec.Command("git", "show", ":"+f).Output()
		result.WriteString(fmt.Sprintf("=== %s ===\n%s\n", f, string(content)))
	}
	return result.String(), nil
}
