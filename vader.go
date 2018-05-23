package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type bindef struct {
	Path  string
	Pyver string
}

func findPythonBins() []bindef {

	var path = os.Getenv("PATH")
	var pythons = make([]bindef, 0)
	for _, s := range strings.Split(path, ":") {
		var programs, _ = ioutil.ReadDir(s)
		for _, p := range programs {
			var pn = p.Name()
			if strings.HasPrefix(pn, "python") && !strings.Contains(pn, "-") {
				var e = bindef{}
				e.Path = s + "/" + pn
				e.Pyver = strings.TrimPrefix(pn, "python")
				pythons = append(pythons, e)
			}
		}
	}

	return pythons

}

type vaderfiledef struct {
	Main  string `yaml:"main"`
	Pyver string `yaml:"pyver"`
}

func parseVaderfile(path string) vaderfiledef {

	// Mostly stolen from https://stackoverflow.com/questions/28682439/

	filename, _ := filepath.Abs(path)
	yamlfile, err := ioutil.ReadFile(filename)

	// MEME
	if err != nil {
		panic(err)
	}

	var out vaderfiledef
	err = yaml.Unmarshal(yamlfile, &out)

	// MEME
	if err != nil {
		panic(err)
	}

	return out

}

func runPython(vf vaderfiledef, bin string) {
	var prog = exec.Command(bin, vf.Main)
	prog.Stdin = os.Stdin
	prog.Stdout = os.Stdout
	prog.Stderr = os.Stderr
	prog.Run()
}

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		println("not enough arguments")
		return
	}

	verb := args[0]

	if verb == "run" {

		var vf = parseVaderfile("./Vaderfile")
		runPython(vf, "python"+vf.Pyver)

	} else if verb == "pull" {

		if len(args) != 4 {
			panic("not enough args")
		}

		pyver := args[1]
		pkgname := args[2]
		pkgver := args[3]

		pkg := pippackage{
			Pipver:  pyver,
			Name:    pkgname,
			Version: pkgver,
		}

		downloadPackage(pkg)
		buildPackage(pkg)

	} else if verb == "diag-lspy" {

		var pys = findPythonBins()
		for _, bd := range pys {
			println("python" + bd.Pyver + " " + bd.Path)
		}

	} else {
		println("bad")
	}

}
