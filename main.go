package main

import (
	"fmt"
	"os/exec"
	"strings"
)

const fromUser = "+493641945284"
const toUser = "+491715761382"

var s status

func textize(input string) (output string) {
	output = strings.ReplaceAll(input, "<strong>", "")
	output = strings.ReplaceAll(output, "</strong>", "")
	return output
}

func sendSignal(format string, a ...interface{}) {
	err := exec.Command("/usr/local/bin/signal-cli", "-u", fromUser, "send", "-m", fmt.Sprintf(format, a...), toUser).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	stadtJena("/home/ber/bin/.coronaStadt")
	otzBlog("/home/ber/bin/.coronaOtz")
}
