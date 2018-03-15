package main

import "os"
import "os/exec"
import "io/ioutil"
import "strings"
import "fmt"

func find_python_bins() []string {

	var path = os.Getenv("PATH")
	var pythons = make([]string, 0)
	for _, s := range strings.Split(path, ":") {
		var programs, _ = ioutil.ReadDir(s)
		for _, p := range programs {
			var pn = p.Name();
			if strings.HasPrefix(pn, "python") && !strings.HasSuffix(pn, "-config") {
				pythons = append(pythons, s + "/" + pn)
			}
		}
	}

	return pythons

}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("error: not enough arguments\n");
		os.Exit(1)
	}

	var target = os.Args[1]
	var prog = exec.Command("/usr/bin/python", target);
	prog.Stdin = os.Stdin
	prog.Stdout = os.Stdout
	prog.Stderr = os.Stderr
	prog.Run()

}
