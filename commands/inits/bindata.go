// Code generated by go-bindata.
// sources:
// data/aliases.sh
// data/git_wrapper.sh
// data/status_shortcuts.sh
// DO NOT EDIT!

package inits

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dataAliasesSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xc9\xc1\x0d\x83\x21\x08\x06\xd0\xbb\x53\x78\x63\x0a\x67\x69\x08\x88\x35\xb5\xb1\xf1\xc3\xfd\x7b\xe0\x0f\xd7\xf7\x78\x4d\x46\x1d\x68\x04\xf9\xfe\xae\xd9\x0b\xce\x7e\x41\xe5\x19\x6e\x34\xa6\x57\x56\x4d\xd2\x20\x9d\x66\x69\x2b\x6c\xed\x91\x24\x3b\x4c\xde\x5d\x3e\xfb\x7a\xc6\x41\xc4\xe9\xe8\x4e\xe5\x1f\x00\x00\xff\xff\xb3\xf6\xf4\x39\x83\x00\x00\x00")

func dataAliasesShBytes() ([]byte, error) {
	return bindataRead(
		_dataAliasesSh,
		"data/aliases.sh",
	)
}

func dataAliasesSh() (*asset, error) {
	bytes, err := dataAliasesShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/aliases.sh", size: 131, mode: os.FileMode(420), modTime: time.Unix(1558658378, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataGit_wrapperSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8f\x41\x6b\x14\x41\x10\x85\xef\xfd\x2b\x1e\x93\xc1\xec\x8a\x4b\x88\xd7\xc5\x20\x44\x22\x1e\x02\xa2\x06\x2f\xc2\x52\x33\x53\xbd\x53\xd8\xd3\x3d\x74\x57\xcf\x26\x38\xfe\x77\xe9\xc9\xc6\x43\x70\x51\xf0\xd2\x87\xaa\xf7\xbe\xaf\xfa\x0c\x9f\x78\x08\x13\x83\xfc\x03\xf8\x5e\x92\x8a\xdf\x63\x2f\x0a\x72\x42\x09\x21\xc2\x66\xdf\xaa\x04\x6f\xb2\x7f\x9c\x95\xed\x15\x2e\x3a\x9e\x2e\x7c\x76\x0e\xaf\xaf\x5e\x5c\x9a\xec\x13\x2b\x36\xf6\xcf\x5b\x73\x86\xbb\xc4\xd0\x9e\x61\xcb\x70\x24\xed\xa1\x61\x09\x6b\x00\x4d\x41\x3a\x88\xb7\xe2\x45\x19\x2e\x84\x11\x07\xd1\x7e\xd9\xff\xf6\xf3\xfd\x18\xa2\xe2\xf3\xf5\xed\xc7\xbb\x9b\x9b\xdd\xfb\x0f\x5f\x76\xd7\xb7\xef\xde\x54\xf5\xea\xdb\xa1\x97\x76\x49\xaf\xab\xe2\xfa\x1a\x69\x5c\xba\x0b\xa4\x58\xcf\xfb\xdc\x9c\x97\x51\x9f\x1b\x1c\x22\x8d\x23\xc7\x57\x10\x0b\xf1\x49\xc9\x39\xee\x8c\x58\xe8\xc3\xc8\x28\x89\xe7\x1f\xd8\x16\x88\xc7\xa9\x0b\xfa\xdc\x54\x5b\x58\x31\xe6\xe9\xd8\xa2\x5a\xad\xf1\xc3\x00\x2d\x25\x46\x7d\x09\xf1\x06\x00\xda\x30\x0c\xa2\x73\xe3\x68\xe0\xd9\x85\xfd\x1c\xb9\xa1\xc4\xf3\xc0\x71\xcf\xeb\x25\x02\xf0\x44\x0e\x55\xbd\x4a\xed\x30\x66\x6b\x8b\x98\x7c\x87\xcd\x06\x55\xfd\xcc\x5e\xa1\xaa\xdf\x56\xeb\x6a\xbb\x7d\xc4\xf7\xdc\x7e\x0f\x59\xe7\x4e\xac\x9d\xe3\x30\x47\x4e\xac\xe5\xd5\x10\xff\xce\x8f\xec\x48\x65\xe2\x7f\x51\x51\xd7\xfd\xc7\xbd\xc7\xe6\xb1\xb1\x4b\x4a\x9a\xd3\x91\xfc\xf2\x89\x7b\xa2\xbd\xc4\x38\x51\x6b\x7e\x9a\x5f\x01\x00\x00\xff\xff\x17\xa5\xb3\x5b\xc2\x02\x00\x00")

func dataGit_wrapperShBytes() ([]byte, error) {
	return bindataRead(
		_dataGit_wrapperSh,
		"data/git_wrapper.sh",
	)
}

func dataGit_wrapperSh() (*asset, error) {
	bytes, err := dataGit_wrapperShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/git_wrapper.sh", size: 706, mode: os.FileMode(420), modTime: time.Unix(1558658378, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataStatus_shortcutsSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x94\x41\x6f\xfb\x44\x10\xc5\xef\xfe\x14\x0f\xff\x2d\x35\x51\x12\x95\x72\xab\x82\x05\x12\x4d\x45\x2f\x14\x25\x85\x03\xb4\x4a\x37\xf6\xb8\x5e\xb1\xd9\xb5\x76\xc6\x71\x21\xe4\xbb\xa3\x5d\xc7\x49\x29\x54\xbd\x25\xda\xd9\x99\xf7\x7e\xf3\xd6\x5c\x6c\x9b\xb6\xaa\xd6\x2c\x4a\x5a\x1e\x8d\xb1\x4f\x00\xe3\x0a\x65\x30\x1c\x91\xdd\xad\x8b\x5a\xf9\x3c\xa5\x34\x49\x80\x2f\x58\x58\x6e\x3d\x81\xeb\xce\xf9\x92\x1b\xa3\x05\x9a\xe1\x2c\x2a\xe7\xf1\x17\xd7\x09\xa0\x2b\xfc\x8e\x99\x45\x9a\xfd\xb6\xfa\x71\xfd\xeb\x62\xb9\xba\xbb\xff\x29\xc5\xd3\x1c\x52\x93\x05\x93\xb8\x46\xde\x76\x98\xa3\xd2\xf3\xbe\xfd\xb2\xb5\xc3\x70\xf4\xba\xa6\x60\x71\x9e\xe0\x5a\x69\x5a\x89\x45\xa3\xe7\xa8\xf2\x19\x96\xa8\x64\x88\xc3\x86\x82\x04\x2d\x0c\xd7\x59\x18\x6d\x09\x4e\x6a\xf2\x9d\x66\x02\xbd\x6a\x41\xe1\x4a\x0a\x4a\xb9\x53\xc6\xb8\x8e\xca\xaf\xc6\x27\xb7\xc5\xb6\x5c\x9f\xda\x9f\xff\xe4\x69\x36\xba\x6c\xd9\x5f\x6e\xb4\xbd\x24\xbb\x7b\x27\x0c\xb3\x59\xa5\x0d\x19\xcd\x82\xec\xfb\xf1\x91\x8f\xae\x82\x4b\x4f\xe8\x14\x43\x59\x90\xf7\xce\x4f\x7b\x0d\x8d\xa7\xad\x92\xd6\x93\xf9\x73\x0a\x65\x4b\x34\x8a\x19\xca\x38\xfb\x12\x2e\x9d\x85\xf6\x2e\x57\x0f\x37\xf7\xbf\x3c\xc4\x46\x27\xd5\xd8\xb4\x02\xeb\x04\xab\x87\x9b\xc5\x72\x39\x05\x3b\xb4\x4c\x1e\x5c\xbb\xd6\x94\x60\xd1\xc6\x80\x89\xfa\xc1\xd8\xf2\xcb\xd9\x27\x71\x9e\x7d\x37\x2c\x28\x23\xc6\xcc\x12\xbe\x1e\x16\x93\x00\x80\x27\x69\xbd\x0d\x87\x09\x50\xe9\xde\xd3\x2d\x49\x51\x23\x1a\x75\x15\x82\x69\xc6\xa8\xf2\x6e\x8b\x4a\x7b\x96\x23\xef\x0a\x5c\x78\xdd\xc8\x71\x53\xe3\xd8\xc0\x10\x07\x8c\x54\xd4\x0e\x69\x76\x66\x9b\xe2\x6f\xd4\xa4\xca\x90\x93\xab\x81\xdd\xe2\xb5\x71\x5e\x60\xdb\xed\x86\x3c\x95\x08\xcc\x77\xca\x6b\xb5\x09\x13\x43\xc0\x48\x15\x75\xec\x9a\xe0\x94\xd1\xc2\x90\xf2\xeb\x9d\xf2\x41\xf2\xdd\xed\x2a\xcf\x2e\x1e\xe5\xe2\x6c\x3a\xbf\x0a\x4a\x9c\x8f\xf7\xa0\x2d\xb2\x28\x6b\x8e\xd2\x45\xcb\xd4\x4f\xcd\xde\x67\x3e\xa3\x3c\x8d\xa5\x69\x2c\x33\x24\xa0\xc9\x24\x01\x4a\x67\xe9\x34\x0a\x8f\xf2\x68\x2f\x7a\xfd\x3f\x7b\x6d\x65\x48\x47\x0f\x28\xa2\x91\xce\xc1\xd9\x4e\xf9\x32\x40\xf9\x5f\x16\xa2\xb4\x09\x2c\x26\xdf\x1c\xdf\x01\x31\x49\x78\x4f\x01\x82\xf6\xce\x6e\xc9\x4a\x48\x7a\x49\x95\x6a\x8d\x7c\xfe\xcc\x5a\xfb\xd1\x43\x3b\x24\x49\xf2\x05\x3f\x04\x6c\x1f\xb0\x4e\xfe\xcb\xf6\xf3\xcf\xc3\x70\x18\x43\x13\x78\x8f\x46\xd0\xf9\xd5\x1c\xfa\xdb\xfc\xfa\xfa\x7a\x0e\x3d\x99\x60\x3c\x3e\x71\x3f\xae\xc7\xee\x42\xff\xb5\xce\xb3\xfd\xfb\xb6\x87\x6c\xaf\x0f\xb1\x36\x78\x8d\x66\xb3\xfd\xe9\xc2\x01\x4f\x6f\x93\x8b\xde\xf1\xbf\x2a\xfa\xfd\x1a\xa6\x63\xc5\xc6\x93\xfa\x23\xfe\xae\xf4\xb0\xc8\x43\xf2\x4f\x00\x00\x00\xff\xff\xfa\x9a\x5d\x4e\x08\x05\x00\x00")

func dataStatus_shortcutsShBytes() ([]byte, error) {
	return bindataRead(
		_dataStatus_shortcutsSh,
		"data/status_shortcuts.sh",
	)
}

func dataStatus_shortcutsSh() (*asset, error) {
	bytes, err := dataStatus_shortcutsShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/status_shortcuts.sh", size: 1288, mode: os.FileMode(420), modTime: time.Unix(1558658378, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"data/aliases.sh": dataAliasesSh,
	"data/git_wrapper.sh": dataGit_wrapperSh,
	"data/status_shortcuts.sh": dataStatus_shortcutsSh,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"data": &bintree{nil, map[string]*bintree{
		"aliases.sh": &bintree{dataAliasesSh, map[string]*bintree{}},
		"git_wrapper.sh": &bintree{dataGit_wrapperSh, map[string]*bintree{}},
		"status_shortcuts.sh": &bintree{dataStatus_shortcutsSh, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

