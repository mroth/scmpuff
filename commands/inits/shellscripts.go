package inits

// the below go:generate directive will automatically generate a bindata.go file
// which wraps the contents of the data directory so we can include text files
// in our distributed binary directly. (boy, Go can be quite annoying
// sometimes!)

//go:generate go-bindata -o bindata.go -ignore=README* -pkg=inits data

func statusShortcutsString() string {
	data, err := Asset("data/status_shortcuts.sh")
	if err != nil {
		// Asset was not found.
	}
	return string(data)
}
