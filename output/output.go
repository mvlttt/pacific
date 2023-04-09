package output

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Out(f *os.File, results chan string, blacklist []string) {
	for result := range results {
		result = strings.TrimSpace(result)
		if result == "" {
			continue
		}
		u, err := url.Parse(result)
		if err != nil {
			continue
		}

		// If the blacklist is empty, "strings.Split" returns a slice with one element. Therefore, I also check the length of the first element.
		if len(blacklist) > 0 && len(blacklist[0]) > 0 {
			if filepath.Ext(u.Path) != "" && contains(blacklist, filepath.Ext(u.Path)[1:]) {
				continue
			}
		}
		writer := bufio.NewWriter(f)
		_, _ = writer.WriteString(result + "\n")
		writer.Flush()
	}
}

func Err(msg string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s %s", red("[Error]"), msg+"\n")
}

func Info(msg string) {
	red := color.New(color.FgBlue).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s %s", red("[Info]"), msg+"\n")
}

func contains(slice []string, ext string) bool {
	for _, e := range slice {
		if e == ext {
			return true
		}
	}
	return false
}
