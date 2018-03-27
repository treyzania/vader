package main

import "os"
import "os/exec"
import "os/user"
import "io/ioutil"
import "path"
import "path/filepath"

import "strings"

import "gopkg.in/yaml.v2"

const pip2path = "/usr/bin/pip"
const pip3path = "/usr/bin/pip3"

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

type pippackage struct {
	Pipver string
	Name string
	Version string
}

func download_package(pkg pippackage) {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	ppath := path.Join(user.HomeDir, ".vader", "repo", pkg.Pipver, pkg.Name, pkg.Version)
	err = os.MkdirAll(ppath, os.ModeDir | 0755)
	if err != nil {
		panic(err)
	}

	// Figure out what we should pass to pip as the version string.
	var dlstr = pkg.Name
	if len(pkg.Version) > 0 {
		dlstr = dlstr + "==" + pkg.Version
	}

	dlcmd := exec.Command(pkg.Pipver, "download", "--no-deps", dlstr)
	dlcmd.Dir = ppath
	dlcmd.Stdout = os.Stdout
	dlcmd.Stderr = os.Stderr
	dlcmd.Run()

}

func run_python(vf vaderfiledef, bin bindef) {
	var prog = exec.Command(bin.Path, vf.Main);
	prog.Stdin = os.Stdin
	prog.Stdout = os.Stdout
	prog.Stderr = os.Stderr
	prog.Run()
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

	run_python(vf, bin)

}
