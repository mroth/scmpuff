package inits

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _data_aliases_sh = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x3c\xc9\x41\x0a\xc4\x20\x0c\x05\xd0\xfd\x9c\xc2\x9d\xa7\xf0\x2c\x43\x48\x8c\x95\x5a\x2c\xfe\x78\xff\x16\x52\xb2\x7d\x8f\x46\x27\xa4\x86\x92\xc1\xd7\xbd\x55\xff\x30\xb2\x8d\xfc\xfb\x86\x4a\x6e\xdd\x12\x89\x04\x89\x93\x74\xd5\xb0\xe1\x36\x66\x0b\xe2\xe9\xc6\x47\xe5\x73\x6e\x8b\x58\xf0\x58\x15\xf5\xd5\x27\x00\x00\xff\xff\xb3\xf6\xf4\x39\x83\x00\x00\x00")

func data_aliases_sh_bytes() ([]byte, error) {
	return bindata_read(
		_data_aliases_sh,
		"data/aliases.sh",
	)
}

func data_aliases_sh() (*asset, error) {
	bytes, err := data_aliases_sh_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "data/aliases.sh", size: 131, mode: os.FileMode(420), modTime: time.Unix(1426541666, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _data_git_wrapper_sh = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x92\x41\x8b\xdb\x30\x10\x85\xef\xfa\x15\x83\x13\x9a\xb8\x34\x84\xf4\x6a\x1a\x0a\x29\x29\x3d\x04\x4a\xdb\xd0\x4b\x21\x28\xf6\xc8\x1e\x2a\x4b\x46\x92\x9d\x84\xf5\xfe\xf7\x1d\x39\xc9\x1e\xc2\x86\xec\x49\xcc\xbc\x37\xdf\xd3\xd8\x1a\xc1\x2f\xac\x6d\x87\x20\xcd\x09\xf0\x48\x3e\x90\x29\xa1\xa4\x00\x52\x93\xf4\x60\x1d\xa8\xd6\xe4\x81\xac\x11\xad\x39\xf7\xa2\xba\x84\x79\x81\xdd\xdc\xb4\x5a\xc3\xe7\xe5\x87\x05\x8b\x1e\x03\xcc\xd4\xdb\xaa\x18\xc1\xd6\x23\x84\x0a\x19\xc7\xcd\x46\x86\x0a\x82\x1d\xcc\x7c\xc8\xce\x52\x01\x64\x14\x19\x0a\x08\xda\xda\x06\x0e\xc4\x96\xa8\xbf\xe6\xe3\xb1\xb1\x2e\xc0\xef\xd5\xe6\xe7\x76\xbd\xde\x7d\xff\xf1\x67\xb7\xda\x7c\xfb\x92\x8c\xa7\xff\x0e\x15\xe5\x83\x3b\x4d\x62\xd6\x5f\x27\x9b\x61\x76\x80\xc4\xd4\x49\xd5\xee\x27\xb1\xc5\x27\x1c\x58\x6e\xd0\x7d\x02\x52\x1c\xea\x83\xd4\x1a\x0b\xc1\x45\x38\x35\x08\xd1\x71\xbb\x40\x16\x21\x06\xee\xdd\x80\x47\x92\x0c\x14\x09\x71\xbd\x6c\x8c\x9a\xa6\xf0\x24\x00\x72\xc9\x9b\x8f\x17\x9c\xc4\x05\x97\xb6\xae\x29\xf4\x7b\x2d\x6b\xec\xb5\x2d\x7b\x87\x7b\x76\xf4\x35\xba\x12\xd3\xc1\x02\x80\x9d\xd4\xc0\x8b\xf9\xbc\x6e\x5a\xa5\x62\xb0\x34\x05\xcc\x66\xdc\xbc\x49\x4f\xb8\xf5\x35\x49\x93\x2c\x3b\xe3\x2b\xcc\xff\xdb\x36\xf4\x05\x29\xd5\xbb\x9a\xf1\xfc\x63\x1e\x72\x1d\x6a\x19\x88\xdf\xc1\x3b\x22\x64\x51\x3c\xe0\xdd\x25\x5c\xc6\x2e\xf6\x1d\x7f\xfa\xd0\xfa\x0b\xf6\xe3\x15\x7a\x67\x7a\xb0\xa1\x97\xb9\x78\x16\x2f\x01\x00\x00\xff\xff\xc6\xff\x30\xc2\xb7\x02\x00\x00")

func data_git_wrapper_sh_bytes() ([]byte, error) {
	return bindata_read(
		_data_git_wrapper_sh,
		"data/git_wrapper.sh",
	)
}

func data_git_wrapper_sh() (*asset, error) {
	bytes, err := data_git_wrapper_sh_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "data/git_wrapper.sh", size: 695, mode: os.FileMode(420), modTime: time.Unix(1426541666, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _data_status_shortcuts_sh = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x94\x41\x6f\xd3\x4e\x10\xc5\xef\xfe\x14\xef\xef\x5a\x6a\xa2\xb4\xea\xbf\xdc\xaa\x60\x81\x04\xa9\xe8\x85\xa2\xa4\x70\x80\x56\xe9\xc6\x1e\xd7\x2b\xd6\xbb\xd6\xee\xba\x2e\x84\x7c\x77\x66\xed\x38\x29\x85\xaa\x27\xdb\x99\xd9\x99\xf7\x7e\x33\x1b\x97\x55\x75\x53\x14\x4b\xe7\x85\x6f\xdc\x68\x8c\x75\x04\x28\x93\x09\x85\x21\x44\xfa\x7e\x99\x95\xc2\xa6\x31\xc5\x11\x47\x0f\x30\xd3\xae\xb1\x04\x57\xb6\xc6\xe6\xae\x56\xd2\x43\x3a\x18\x8d\xc2\x58\xfc\x74\x25\x27\xc9\x02\xdf\x70\xac\x11\x27\x5f\x17\x1f\x96\x5f\x66\xf3\xc5\xc5\xe5\xc7\x18\x37\x53\xf8\x92\x34\x1c\x79\x53\xfb\xc7\x15\xa6\x28\xe4\xb4\x2f\x3f\x6f\xf4\xd0\x1c\xbd\xae\x23\x7e\x1a\xee\x68\x1a\x5f\x37\xbe\x4b\x1a\xdd\x76\x2a\x6f\xa1\x89\x72\x07\x6f\xb0\xa2\x20\x41\x7a\x56\xd2\x6a\x28\xa9\xf9\x9b\x9b\xd9\x56\x3a\x02\x3d\xb0\xc8\xcc\xe4\x14\x94\xba\x56\x28\x65\x5a\xca\xff\x1b\xef\xdc\x66\x55\xbe\xdc\x95\xdf\x7f\xa4\x71\x32\x3a\x69\x9c\x3d\x59\x49\x7d\xc2\x24\x9e\x08\xc3\xf1\x71\x21\x15\x29\xe9\x3c\x92\xb7\xe3\x2d\x1f\x36\x1f\x1a\x13\x5a\xe1\x20\x34\xc8\x5a\x63\x8f\x7a\x0d\xb5\xa5\x8a\x4f\x5a\x52\x3f\x8e\x38\x96\xa3\x16\x8e\x93\x94\xd1\x77\xe1\xd0\x5e\x68\xef\x72\x71\xf5\xfe\xf2\xf3\x55\x57\x68\xa7\x1a\xab\xc6\x43\x1b\x0f\x0e\xce\xe6\x73\x86\x63\xd0\x38\xb2\x8c\xd3\x34\x2a\x67\x69\x52\xf1\xf8\x88\xfa\xc6\xa8\xdc\xdd\xde\x27\xb9\x34\x79\x33\x0c\x28\x21\x76\xc0\x9c\xfe\x1f\x06\xc3\x01\xc0\x12\xeb\xd3\x21\xc8\x9f\x85\xec\x3d\x9d\x93\xcf\x4a\x74\x46\x4d\x81\x60\xda\x61\x54\x58\x53\xf1\xbb\xe5\x1f\x7b\xde\xcc\x25\xb3\x92\x27\xdb\xd3\x1b\x77\x05\x38\x35\x60\xa4\xac\x34\xbc\x0f\x7b\xb6\x31\x7e\xa1\x24\x91\x87\x3d\x39\x1d\xd8\xcd\x1e\x6a\x63\xd9\x5e\x53\xad\x98\x60\x8e\xc0\xfc\x5e\x58\x29\x56\xa1\x63\x58\x30\x12\x2c\x24\x54\xe5\xfc\x61\x47\x33\x45\xc2\x2e\x39\x2f\x48\xbe\x38\x5f\xa4\xc9\xe1\xb5\x3f\xdc\x9b\x4e\x4f\x83\x12\x3e\x1c\xce\x41\xb2\xb9\x4e\xd6\x14\xb9\xe9\x2c\x53\xdf\x35\x79\xba\xf3\x09\xb1\xf2\x90\x1a\x77\x69\x8a\x3c\x68\x32\xe1\xf7\xdc\x68\xda\xb5\xc2\xb5\xbf\xd6\x87\xbd\xfe\x4f\x56\x6a\x3f\x6c\x47\x0f\xa8\x43\xe3\x5b\xc3\xeb\xd9\x0a\x9b\x07\x28\xff\x64\xe1\x85\x54\x81\xc5\xe4\xd5\xf6\x1e\x10\x5f\x93\x70\x9f\x02\x04\x69\x8d\xae\x88\x4b\xf3\xa6\xe7\x54\x88\x46\xf9\x97\xaf\x59\xa3\x9f\xbb\x68\x9b\x28\x8a\x0e\xf0\x2e\x60\x7b\x86\x75\xf4\x37\xdb\x97\xff\x1e\x86\x60\xb7\x34\x81\xf7\x68\x04\x99\x9e\x4e\x21\x5f\xa7\x67\x67\x67\xfc\x9c\x4c\x30\x1e\xef\xb8\x6f\xc7\xc3\x25\xb8\xfe\x52\xa6\xc9\xfa\x69\xd9\x4d\xb2\x96\x9b\x2e\x37\x78\xed\xcc\x26\xeb\xdd\x81\x0d\x6e\x1e\x6f\x2e\x7a\xc7\x7f\x64\xf4\xf3\x55\x8e\xb6\x19\x2b\x4b\xe2\x7b\xf7\xce\xab\xbd\x1d\xe4\x26\xfa\x1d\x00\x00\xff\xff\xfa\x9a\x5d\x4e\x08\x05\x00\x00")

func data_status_shortcuts_sh_bytes() ([]byte, error) {
	return bindata_read(
		_data_status_shortcuts_sh,
		"data/status_shortcuts.sh",
	)
}

func data_status_shortcuts_sh() (*asset, error) {
	bytes, err := data_status_shortcuts_sh_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "data/status_shortcuts.sh", size: 1288, mode: os.FileMode(420), modTime: time.Unix(1426541666, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	if (err != nil) {
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
	"data/aliases.sh": data_aliases_sh,
	"data/git_wrapper.sh": data_git_wrapper_sh,
	"data/status_shortcuts.sh": data_status_shortcuts_sh,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"data": &_bintree_t{nil, map[string]*_bintree_t{
		"aliases.sh": &_bintree_t{data_aliases_sh, map[string]*_bintree_t{
		}},
		"git_wrapper.sh": &_bintree_t{data_git_wrapper_sh, map[string]*_bintree_t{
		}},
		"status_shortcuts.sh": &_bintree_t{data_status_shortcuts_sh, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
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

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

