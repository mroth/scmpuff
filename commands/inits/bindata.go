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

var _data_aliases_sh = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x3c\xc9\xc1\x0d\x80\x20\x0c\x05\xd0\xbb\x53\x70\x63\x0a\x66\x31\x4d\x4b\x91\xa8\xc1\xf0\xcb\xfe\x9a\xd4\xf4\xfa\x1e\x5d\x9d\x90\x1a\x4a\x06\xdf\xcf\x52\xdd\x61\x64\x0b\x79\xfb\x87\x4a\x6e\xdd\x12\x89\x04\x89\x93\x74\xd5\xb0\x09\xc7\x59\x51\x2d\x94\x87\x2b\x1f\x95\xcf\xb1\xbe\x78\x03\x00\x00\xff\xff\x39\xab\x99\xc7\x70\x00\x00\x00")

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

	info := bindata_file_info{name: "data/aliases.sh", size: 112, mode: os.FileMode(420), modTime: time.Unix(1425406894, 0)}
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

	info := bindata_file_info{name: "data/git_wrapper.sh", size: 695, mode: os.FileMode(420), modTime: time.Unix(1425406894, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _data_status_shortcuts_sh = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x94\x51\x4f\xdb\x4a\x10\x85\xdf\xfd\x2b\xce\x35\x96\x48\x14\x10\x97\xfb\x86\x72\xad\x56\x6a\x83\xca\x4b\xa9\x12\xda\x87\x16\x14\x36\xf6\x18\xaf\xba\xde\xb5\x76\xd7\x98\x36\xe4\xbf\x77\xd6\x8e\x13\x4a\x8b\x78\xb2\x9d\x99\x9d\x39\xe7\x9b\xd9\xb8\xac\xaa\x9b\xa2\x58\x3a\x2f\x7c\xe3\x46\x63\xac\x23\x40\x99\x4c\x28\x0c\x21\xd2\xf7\xcb\xac\x14\x36\x8d\x29\x8e\x38\x7a\x80\x99\x76\x8d\x25\xb8\xb2\x35\x36\x77\xb5\x92\x1e\xd2\xc1\x68\x14\xc6\xe2\xa7\x2b\x39\x49\x16\xf8\x86\x63\x8d\x38\xf9\xba\xf8\xb0\xfc\x32\x9b\x2f\x2e\x2e\x3f\xc6\xb8\x99\xc2\x97\xa4\xe1\xc8\x9b\xda\x3f\xad\x30\x45\x21\xa7\x7d\xf9\x79\xa3\x87\xe6\xe8\x75\x1d\xf1\xd3\x70\x47\xd3\xf8\xba\xf1\x5d\xd2\xe8\xb6\x53\x79\x0b\x4d\x94\x3b\x78\x83\x15\x05\x09\xd2\xb3\x92\x56\x43\x49\xcd\xdf\xdc\xcc\xb6\xd2\x11\xe8\x81\x45\x66\x26\xa7\xa0\xd4\xb5\x42\x29\xd3\x52\xfe\xcf\x78\xe7\x36\xab\xf2\xe5\xae\xfc\xfe\x23\x8d\x93\xd1\x49\xe3\xec\xc9\x4a\xea\x13\x26\xf1\x4c\x18\x8e\x8f\x0b\xa9\x48\x49\xe7\x91\xbc\x1d\x6f\xf9\xb0\xf9\xd0\x98\xd0\x0a\x07\xa1\x41\xd6\x1a\x7b\xd4\x6b\xa8\x2d\x55\x7c\xd2\x92\xfa\x71\xc4\xb1\x1c\xb5\x70\x9c\xa4\x8c\xbe\x0b\x87\xf6\x42\x7b\x97\x8b\xab\xf7\x97\x9f\xaf\xba\x42\x3b\xd5\x58\x35\x1e\xda\x78\x70\x70\x36\x9f\x33\x1c\x83\xc6\x91\x65\x9c\xa6\x51\x39\x4b\x93\x8a\xc7\x47\xd4\x37\x46\xe5\xee\xf6\x3e\xc9\xa5\xc9\x9b\x61\x40\x09\xb1\x03\xe6\xf4\xef\x30\x18\x0e\x00\x96\x58\x9f\x0e\x41\xfe\x2c\x64\xef\xe9\x9c\x7c\x56\xa2\x33\x6a\x0a\x04\xd3\x0e\xa3\xc2\x9a\x8a\xdf\x2d\xff\xd8\xf3\x66\x2e\x99\x95\x3c\xd9\x9e\xde\xb8\x2b\xc0\xa9\x01\x23\x65\xa5\xe1\x7d\xd8\xb3\x8d\xf1\x88\x92\x44\x1e\xf6\xe4\x74\x60\x37\x7b\xa8\x8d\x65\x7b\x4d\xb5\x62\x82\x39\x02\xf3\x7b\x61\xa5\x58\x85\x8e\x61\xc1\x48\xb0\x90\x50\x95\xf3\x87\x1d\xcd\x14\x09\xbb\xe4\xbc\x20\xf9\xe2\x7c\x91\xc6\x8f\xf1\xde\x72\x7a\x1a\x74\xf0\xd1\x70\x0a\x92\xad\x75\xa2\xa6\xc8\x4d\x67\x98\xfa\x9e\xc9\xf3\x8d\x4f\x88\x75\x87\xd4\xb8\x4b\x53\xe4\x41\x93\x09\xbf\xe7\x46\xd3\xb6\x51\x72\x88\x6b\x7f\xad\x0f\x7b\xf5\x9f\xac\xd4\x7e\xd8\x8d\x1e\x4f\x07\xc6\xb7\x86\x97\xb3\x15\x36\x0f\x48\xfe\x4a\xc2\x0b\xa9\x02\x89\xc9\x7f\xdb\x5b\x40\x7c\x49\xc2\x6d\x0a\x08\xa4\x35\xba\x22\x2e\xcd\x7b\x9e\x53\x21\x1a\xe5\x5f\xbf\x64\x8d\x7e\xe9\x9a\x6d\xa2\x28\x3a\xc0\xbb\x00\xed\x05\xd2\xd1\x9f\x64\x5f\xff\x73\x18\x82\xdd\xca\x04\xde\xa3\x11\x64\x7a\x3a\x85\xfc\x3f\x3d\x3b\x3b\xe3\xe7\x64\x82\xf1\x78\xc7\x7d\x3b\x1e\x2e\xc1\xf5\x97\x32\x4d\xd6\xcf\xcb\x6e\x92\xb5\xdc\x74\xb9\xc1\x6b\x67\x36\x59\xef\x0e\x6c\x70\xf3\x74\x6f\xd1\x3b\xfe\x2d\xa3\x9f\xaf\x72\xb4\xcd\x58\x59\x12\xdf\xbb\x77\x5e\xec\xed\x20\x37\xd1\xaf\x00\x00\x00\xff\xff\x29\x2b\x41\xe0\x06\x05\x00\x00")

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

	info := bindata_file_info{name: "data/status_shortcuts.sh", size: 1286, mode: os.FileMode(420), modTime: time.Unix(1425406894, 0)}
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

