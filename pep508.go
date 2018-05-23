package main

type requirements struct {
	Name      string
	Selectors []reqselector
	URL       *string
}

type reqselector struct {
	Kind    string
	Version string
}

// TODO Parsing these ugly things.
