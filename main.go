package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var fromUser string
var toUser string
var saveFile string

func textize(input string) (output string) {
	output = strings.ReplaceAll(input, "<strong>", "")
	output = strings.ReplaceAll(output, "</strong>", "")
	return output
}

func sendSignal(format string, a ...interface{}) {
	err := exec.Command("/usr/local/bin/signal-cli", "-u", fromUser, "send", "-m", fmt.Sprintf(format, a...), toUser).Run()
	// _, err := fmt.Println("/usr/local/bin/signal-cli", "-u", fromUser, "send", "-m", fmt.Sprintf(format, a...), toUser)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	setConfig()
	s := load(saveFile)
	stadtJena(&s.Stadt)
	otzBlogThueringen(&s.OtzThueringen)
	otzBlogWeltweit(&s.OtzWeltweit)
	s.save()
}
