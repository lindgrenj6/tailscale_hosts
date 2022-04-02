package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
	Update my /etc/hosts with all of my tailscale nodes in the form:

	100.x.y.z hostname hostname.local
*/
func main() {
	if os.Getuid() != 0 {
		fmt.Println("need root")
		os.Exit(1)
	}

	raw := must(exec.Command("tailscale", "status").Output())
	lines := strings.Split(string(raw), "\n")

	for i := range lines {
		if !strings.HasPrefix(lines[i], "100.") {
			continue
		}

		entry := strings.Fields(lines[i])[:2]
		entry = append(entry, entry[len(entry)-1]+".local")

		lines[i] = strings.Join(entry, " ")
	}

	must(exec.Command("sed", "-i", "/100./d", "/etc/hosts").Output())
	f := must(os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644))
	defer f.Close()

	must(f.WriteString(strings.Join(lines, "\n")))

	fmt.Printf("%#v", lines)
}

func must[T any](thing T, err error) T {
	if err != nil {
		panic(err)
	}
	return thing
}
