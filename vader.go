package main

import "os"
import "os/exec"
import "io/ioutil"
import "path/filepath"

import "strings"

import "gopkg.in/yaml.v2"

type bindef struct {
	Path string
	Genre string
}

func find_python_bins() []bindef {

	var path = os.Getenv("PATH")
	var pythons = make([]bindef, 0)
	for _, s := range strings.Split(path, ":") {
		var programs, _ = ioutil.ReadDir(s)
		for _, p := range programs {
			var pn = p.Name();
			if strings.HasPrefix(pn, "python") && !strings.HasSuffix(pn, "-config") {
				var e = bindef {}
				e.Path = s + "/" + pn
				e.Genre = strings.TrimPrefix(pn, "python")
				pythons = append(pythons, e)
			}
		}
	}

	return pythons

}

type vaderfiledef struct {
	Main string `yaml:"main"`
	Pyver string `yaml:"pyver"`
}

func parse_vaderfile(path string) vaderfiledef {

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

func main() {

	var vf = parse_vaderfile("./Vaderfile")
	var bin bindef

	for _, pg := range find_python_bins() {
		if pg.Genre == vf.Pyver {
			bin = pg
			break
		}
	}

	var prog = exec.Command(bin.Path, vf.Main);
	prog.Stdin = os.Stdin
	prog.Stdout = os.Stdout
	prog.Stderr = os.Stderr
	prog.Run()

}
