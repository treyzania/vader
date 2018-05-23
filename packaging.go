package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

type pippackage struct {
	Pipver  string
	Name    string
	Version string
}

func (pkg *pippackage) pkgRepoPath() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(user.HomeDir, ".vader", "repo", "pip"+pkg.Pipver, pkg.Name, pkg.Version)
}

type pkgmeta struct {
	Type string
}

func (pkg *pippackage) pkgRepoMetaPath() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(user.HomeDir, ".vader", "repo", "pip"+pkg.Pipver, pkg.Name, pkg.Version+".meta")
}

func (pkg *pippackage) getMeta() *pkgmeta {
	raw, _ := ioutil.ReadFile(pkg.pkgRepoMetaPath())
	var data pkgmeta
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil
	}
	return &data // shrug?  why are we doing this?
}

func (pkg *pippackage) setMeta(pm pkgmeta) {
	raw, _ := json.Marshal(pm)
	ioutil.WriteFile(pkg.pkgRepoMetaPath(), raw, 0644)
}

func downloadPackage(pkg pippackage) string {

	tempdir, err := ioutil.TempDir("", "vader.tmp.")
	if err != nil {
		panic(err)
	}

	// Figure out what we should pass to pip as the version string.
	var dlstr = pkg.Name
	if len(pkg.Version) > 0 {
		dlstr = dlstr + "==" + pkg.Version
	}

	// Actually download it.
	dlcmd := exec.Command("pip"+pkg.Pipver, "download", "--no-deps", dlstr)
	dlcmd.Dir = tempdir
	dlcmd.Stdout = os.Stdout
	dlcmd.Stderr = os.Stderr
	dlcmd.Run()

	// Find the things we downloaded.
	files, err := ioutil.ReadDir(tempdir)
	if err != nil {
		panic(err)
	}

	// Prepare the metadata before we write it.
	meta := pkgmeta{
		Type: "unknown",
	}

	// And extract the files we need.
	ppath := pkg.pkgRepoPath()
	os.MkdirAll(ppath, 0755)
	for _, f := range files { // This should only run once.

		var extcmd *exec.Cmd
		if strings.HasSuffix(f.Name(), ".tar.gz") {
			extcmd = exec.Command("tar", "-xvzf", f.Name(), "--strip-components=1", "-C", ppath)
			meta.Type = "normal"
		} else if strings.HasSuffix(f.Name(), ".whl") {
			extcmd = exec.Command("unzip", f.Name(), "-d", ppath)
			meta.Type = "wheel"
		} else {
			panic("unsupported package type! (" + f.Name() + ")")
		}

		// Now actually do it.
		extcmd.Dir = tempdir
		extcmd.Stdout = os.Stdout
		extcmd.Stderr = os.Stderr
		extcmd.Run()

		// TODO Sometimes we don't have a setup.py and we're distributed as a ".tar.gz", deal with this in determining the package type.

	}

	os.RemoveAll(tempdir)
	pkg.setMeta(meta)
	return ppath

}

func buildPackage(pkg pippackage) {
	var bcmd = exec.Command("python"+pkg.Pipver, "./setup.py", "build")
	bcmd.Dir = pkg.pkgRepoPath()
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	bcmd.Run()
}
