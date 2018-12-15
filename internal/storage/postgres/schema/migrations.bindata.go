// Code generated by go-bindata.
// sources:
// migrations/1543497171_init.down.sql
// migrations/1543497171_init.up.sql
// DO NOT EDIT!

package schema

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

var __1543497171_initDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x4a\x2d\x4e\x2d\x29\xb6\xe6\xc2\x2a\x59\x5a\x9c\x5a\x04\x97\x8b\x0c\x40\x97\x8a\x2f\x2e\x49\x2c\x29\x2d\xb6\x56\xe0\x02\x04\x00\x00\xff\xff\x43\x7b\xc7\xbb\x5b\x00\x00\x00")

func _1543497171_initDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1543497171_initDownSql,
		"1543497171_init.down.sql",
	)
}

func _1543497171_initDownSql() (*asset, error) {
	bytes, err := _1543497171_initDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1543497171_init.down.sql", size: 91, mode: os.FileMode(420), modTime: time.Unix(1544874101, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1543497171_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\x41\xaf\x12\x31\x14\x85\xd7\xb7\xbf\xe2\xee\x60\x08\x09\x46\xa3\x2e\x58\x55\xb8\x68\xe3\x50\xb0\xd3\x79\x3e\x56\x4d\xa5\x7d\xbe\x46\x60\xc8\xb4\xe8\xfb\xf9\x46\x86\xc1\x21\x62\x30\xea\x76\xe6\x9c\x93\xde\xf3\x9d\xa9\x5a\x2c\x51\xaf\x96\x84\x62\x86\x74\x2f\x0a\x5d\xe0\x21\xfa\xda\xc4\x64\xd3\x21\x8e\xd9\x44\x11\xd7\xd4\x48\x3a\x3f\x90\x17\x48\xb2\x9c\x63\xbf\x27\xe9\x63\x6f\x88\xbd\x3b\x52\x62\x26\x68\xda\xcb\xc6\xec\xec\xe2\x6f\xf2\x63\xb2\x5c\xe8\x6e\x7a\xc4\x3e\x83\xe0\x10\xa0\x20\x25\x78\x8e\x4b\x25\xe6\x5c\xad\xf0\x3d\xad\x86\x0c\xec\x7a\x5d\x1d\x76\xc9\x04\x07\x93\x77\x5c\x61\xff\xc5\xab\x0c\x4b\x29\x3e\x94\x74\x4c\x92\x65\x9e\x0f\x71\x34\xc0\xcf\x9b\xea\x93\xdd\x60\x70\x7e\x97\xc2\x43\xf0\x35\x0e\x46\x8c\x81\xdf\xda\xb0\x41\x00\xb8\xe3\xaa\x09\x78\xf9\x2c\xfb\xe9\x64\xb0\xb7\x31\x7e\xab\x6a\x67\x1e\x6d\x7c\xc4\x8e\xee\xf5\xf3\x0b\xdd\xe9\x56\x00\xe8\x9e\xde\x11\xa4\xea\x8b\xdf\x01\x60\xfb\xce\xae\x1d\x19\x83\xd1\x00\x53\xd8\xfa\x98\xec\x76\x1f\x7f\x3c\xae\x71\x18\xff\xb4\x0f\xb5\x77\xc6\x26\x00\x2d\xe6\x54\x68\x3e\x5f\x76\x83\xd7\xb5\xb7\xe9\x24\xb8\xa2\xc0\x29\xcd\x78\x99\x6b\x9c\x94\x4a\x91\xd4\xe6\x2c\x19\x32\x38\xec\xdd\xdf\x79\x59\x07\x9d\x90\x53\xba\xbf\x86\xce\x1c\xdb\x35\xc1\x3d\xe1\x42\xb6\x34\x8f\xdf\xb2\xf1\x6d\xf3\x57\x5f\x87\x87\xb0\xb6\x29\x54\x3b\xd3\x74\x71\x99\xf4\xab\xe0\xc6\x9e\x6a\x1f\x7d\x3a\x0d\xea\xea\x9c\x5a\x4a\xbf\x81\xcc\x1a\xb6\xc1\x81\x90\x9a\xde\x92\x42\x45\x33\x52\x24\x27\x74\x1e\x6b\x70\xd9\x9f\xe2\xbc\x49\xf3\x1f\x60\xfe\x4f\x96\x4d\x6d\x97\x08\xda\x2a\xdb\xda\xbf\x07\x00\x00\xff\xff\x61\xc1\xc4\xe4\x1d\x04\x00\x00")

func _1543497171_initUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1543497171_initUpSql,
		"1543497171_init.up.sql",
	)
}

func _1543497171_initUpSql() (*asset, error) {
	bytes, err := _1543497171_initUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1543497171_init.up.sql", size: 1053, mode: os.FileMode(420), modTime: time.Unix(1544876145, 0)}
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
	"1543497171_init.down.sql": _1543497171_initDownSql,
	"1543497171_init.up.sql": _1543497171_initUpSql,
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
	"1543497171_init.down.sql": &bintree{_1543497171_initDownSql, map[string]*bintree{}},
	"1543497171_init.up.sql": &bintree{_1543497171_initUpSql, map[string]*bintree{}},
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
