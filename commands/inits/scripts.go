package inits

import "fmt"

// the below go:generate directive will automatically generate a bindata.go file
// which wraps the contents of the data directory so we can include text files
// in our distributed binary directly.

//go:generate go-bindata -o bindata.go -ignore=README* -pkg=inits data

func printScript() {
	fmt.Println(assetString("data/status_shortcuts.sh"))

	if includeAliases {
		fmt.Println(assetString("data/aliases.sh"))
	}

	if wrapGit {
		fmt.Println(assetString("data/git_wrapper.sh"))
	}
}

// returns the string data for an embedded data script
func assetString(file string) string {
	data, err := Asset(file)
	if err != nil {
		// Asset was not found. This should be impossible unless something goes
		// wrong during compilation build process, so panic!
		panic(fmt.Sprintf("Could not find bindata asset file: %v", file))
	}
	return string(data)
}
