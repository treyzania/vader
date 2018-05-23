package main

import (
	"regexp"
)

type requirementdef struct {
	Name      string
	Selectors []reqselector
	URL       *string
}

type reqselector struct {
	Kind    string
	Version string
}

const pkgNameRegex = "^[a-zA-Z0-0\\.\\_\\-]+"
const pkgVerRegex = "^[0-9](\\.[0-9])+"

func parseRequirement(line string) requirementdef {

	// Parse the name of the package.
	var name string
	for i := 1; i < len(line); i++ {

		match, err := regexp.MatchString(line[:i], pkgNameRegex)
		if err != nil {
			panic("regex compilation error")
		}

		if match {
			name = line[:i]
		} else {
			break
		}

	}

	// TODO Figure out how to parse selectors.

	return requirementdef{
		Name:      name,
		Selectors: make([]reqselector, 0),
		URL:       nil,
	}
}

// TODO Parsing these ugly things.
