package main

import "os"
import "os/exec"
import "os/user"
import "io/ioutil"
import "path"
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

type pippackage struct {
	Pipver string
	Name string
	Version string
}

func (pkg *pippackage) pkg_repo_path() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(user.HomeDir, ".vader", "repo", "pip" + pkg.Pipver, pkg.Name, pkg.Version)
}

func download_package(pkg pippackage) string {

	tempdir, err := ioutil.TempDir("", "vaderdl")
	if err != nil {
		panic(err)
	}

	// Figure out what we should pass to pip as the version string.
	var dlstr = pkg.Name
	if len(pkg.Version) > 0 {
		dlstr = dlstr + "==" + pkg.Version
	}

	// Actually download it.
	dlcmd := exec.Command("pip" + pkg.Pipver, "download", "--no-deps", dlstr)
	dlcmd.Dir = tempdir
	dlcmd.Stdout = os.Stdout
	dlcmd.Stderr = os.Stderr
	dlcmd.Run()

	// Find the things we downloaded.
	files, err := ioutil.ReadDir(tempdir)
	if err != nil {
		panic(err)
	}

	// And extract the files we need.
	ppath := pkg.pkg_repo_path()
	os.MkdirAll(ppath, 0755)
	for _, f := range files {

		if strings.HasSuffix(f.Name(), ".tar.gz") {
			extcmd := exec.Command("tar", "-xvzf", f.Name(), "--strip-components=1", "-C", ppath)
			extcmd.Dir = tempdir
			extcmd.Stdout = os.Stdout
			extcmd.Stderr = os.Stderr
			extcmd.Run()
		} else {
			// TODO wheels
		}

	}

	return ppath

}

func build_package(pkg pippackage) {
	var bcmd = exec.Command("python" + pkg.Pipver, "./setup.py", "build")
	bcmd.Dir = pkg.pkg_repo_path()
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	bcmd.Run()
}

func run_python(vf vaderfiledef, bin string) {
	var prog = exec.Command(bin, vf.Main);
	prog.Stdin = os.Stdin
	prog.Stdout = os.Stdout
	prog.Stderr = os.Stderr
	prog.Run()
}

func main() {

	args := os.Args[1:]
	verb := args[0]

	if verb == "run" {

		var vf = parse_vaderfile("./Vaderfile")
		run_python(vf, "python" + vf.Pyver)

	} else if verb == "pull" {

		if len(args) != 4 {
			panic("not enough args")
		}

		pyver := args[1]
		pkgname := args[2]
		pkgver := args[3]

		pkg := pippackage{
			Pipver: pyver,
			Name: pkgname,
			Version: pkgver,
		}

		download_package(pkg)
		build_package(pkg)

	} else {
		println("bad")
	}

}
