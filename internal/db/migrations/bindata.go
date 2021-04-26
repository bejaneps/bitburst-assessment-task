// Code generated by go-bindata. DO NOT EDIT.
//  memcopy: true
//  compress: true
//  decompress: once
//  asset-dir: true
//  restore: true
// sources:
//  internal/db/migrations/1_create_schema_bitburst.down.sql
//  internal/db/migrations/1_create_schema_bitburst.up.sql
//  internal/db/migrations/2_create_table_objects.down.sql
//  internal/db/migrations/2_create_table_objects.up.sql
//  internal/db/migrations/3_create_index_last_seen.down.sql
//  internal/db/migrations/3_create_index_last_seen.up.sql

package migrations

import (
	"bytes"
	"compress/flate"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/tmthrgd/go-bindata/restore"
)

type asset struct {
	name string
	data string
	size int64

	once  sync.Once
	bytes []byte
	err   error
}

func (a *asset) Name() string {
	return a.name
}

func (a *asset) Size() int64 {
	return a.size
}

func (a *asset) Mode() os.FileMode {
	return 0
}

func (a *asset) ModTime() time.Time {
	return time.Time{}
}

func (*asset) IsDir() bool {
	return false
}

func (*asset) Sys() interface{} {
	return nil
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]*asset{
	"1_create_schema_bitburst.down.sql": &asset{
		name: "1_create_schema_bitburst.down.sql",
		data: "" +
			"\x72\x09\xf2\x0f\x50\x08\x76\xf6\x70\xf5\x75\x54\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e" +
			"\x09\x56\x48\xca\x2c\x49\x2a\x2d\x2a\x2e\xb1\x06\x04\x00\x00\xff\xff",
		size: 35,
	},
	"1_create_schema_bitburst.up.sql": &asset{
		name: "1_create_schema_bitburst.up.sql",
		data: "" +
			"\x72\x0e\x72\x75\x0c\x71\x55\x08\x76\xf6\x70\xf5\x75\x54\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0" +
			"\x0c\x0e\x09\x56\x48\xca\x2c\x49\x2a\x2d\x2a\x2e\xb1\x06\x04\x00\x00\xff\xff",
		size: 37,
	},
	"2_create_table_objects.down.sql": &asset{
		name: "2_create_table_objects.down.sql",
		data: "" +
			"\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xca\x2c" +
			"\x49\x2a\x2d\x2a\x2e\xd1\x53\xca\x4f\xca\x4a\x4d\x2e\x29\x56\x52\x70\x76\x0c\x76\x76\x74\x71\xb5" +
			"\x06\x04\x00\x00\xff\xff",
		size: 48,
	},
	"2_create_table_objects.up.sql": &asset{
		name: "2_create_table_objects.up.sql",
		data: "" +
			"\x6c\x8f\xcd\x6a\x02\x31\x18\x45\xd7\x93\xa7\xb8\xb8\xaa\xa0\x7d\x01\x57\x51\x3f\x4b\xe8\x4c\x46" +
			"\x32\xdf\x80\x76\x33\xcc\x4f\x84\x94\x69\x02\x49\x04\x1f\xbf\x74\x28\x42\xa1\xeb\x73\x2e\xdc\x73" +
			"\x30\x24\x99\xc0\x72\x5f\x12\xd4\x09\xba\x66\xd0\x45\x35\xdc\x60\x70\x79\xb8\xc7\x94\x5f\x57\x61" +
			"\xf8\xb4\x63\x4e\x2b\xbc\x88\x22\x74\x6e\x82\xd2\x4c\x6f\x64\x70\x36\xaa\x92\xe6\x8a\x77\xba\x2e" +
			"\x4b\xdd\x96\xe5\x46\x14\xc1\xcf\xce\x5b\xec\xeb\xba\x24\xa9\x9f\x04\x47\x3a\xc9\xb6\x64\xb0\x69" +
			"\x69\x23\x8a\xb9\x4f\xb9\x4b\xd6\x7a\xb0\xaa\xa8\x61\x59\x9d\xf9\xe3\xaf\x79\x68\x8d\x21\xcd\xdd" +
			"\x53\x10\xeb\x9d\x10\xdb\x2d\xc6\x68\xfb\x6c\xd1\x7b\x38\x3f\xd9\x07\x82\xc7\xf2\x6c\x0c\xf3\xfd" +
			"\xcb\xe3\x16\x22\x6e\x7d\xca\x36\x22\xda\x7e\x4a\xe2\xb7\x53\xe9\x23\x5d\x16\xb3\x73\xd3\x03\xb5" +
			"\xfe\xb7\xf2\x87\xaf\x77\xdf\x01\x00\x00\xff\xff",
		size: 283,
	},
	"3_create_index_last_seen.down.sql": &asset{
		name: "3_create_index_last_seen.down.sql",
		data: "" +
			"\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xca\x2c" +
			"\x49\x2a\x2d\x2a\x2e\xd1\xcb\x49\x2c\x2e\x89\x2f\x4e\x4d\xcd\x8b\xcf\x4c\xa9\xb0\x06\x04\x00\x00" +
			"\xff\xff",
		size: 44,
	},
	"3_create_index_last_seen.up.sql": &asset{
		name: "3_create_index_last_seen.up.sql",
		data: "" +
			"\x44\xcb\xbf\xaa\xc2\x30\x14\x07\xe0\xfd\x3e\xc5\x8f\x4e\xd7\xa1\xbe\x80\x93\x68\x84\x2c\x29\xd8" +
			"\x0c\xdd\x4a\xfe\x9c\x42\x24\x9e\x40\xce\x29\xf4\xf1\xdd\x74\xff\xbe\x71\x44\xea\x14\x94\x10\x18" +
			"\x85\x33\x1d\x68\x8c\x1a\x44\x57\x21\x62\xa4\x56\xf7\x37\x63\x6b\x1d\x5b\x10\xa5\x8e\x4c\x95\x94" +
			"\xe4\xef\xf6\x34\x57\x6f\x60\xdd\xdd\x2c\xb0\x0f\xb8\xc9\xc3\x2c\x76\xf6\xf3\xaf\xaf\x25\x1f\x98" +
			"\x1c\x62\xd1\xb8\x77\xd1\xf3\xd0\xe2\x8b\x92\xca\x80\xff\x2f\x3a\x5d\x3e\x01\x00\x00\xff\xff",
		size: 133,
	},
}

// AssetAndInfo loads and returns the asset and asset info for the
// given name. It returns an error if the asset could not be found
// or could not be loaded.
func AssetAndInfo(name string) ([]byte, os.FileInfo, error) {
	a, ok := _bindata[filepath.ToSlash(name)]
	if !ok {
		return nil, nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}

	a.once.Do(func() {
		fr := flate.NewReader(strings.NewReader(a.data))

		var buf bytes.Buffer
		if _, a.err = io.Copy(&buf, fr); a.err != nil {
			return
		}

		if a.err = fr.Close(); a.err == nil {
			a.bytes = buf.Bytes()
		}
	})
	if a.err != nil {
		return nil, nil, &os.PathError{Op: "read", Path: name, Err: a.err}
	}

	return a.bytes, a, nil
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	a, ok := _bindata[filepath.ToSlash(name)]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}

	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	data, _, err := AssetAndInfo(name)
	return data, err
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

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}

	return names
}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	return restore.Asset(dir, name, AssetAndInfo)
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	return restore.Assets(dir, name, AssetDir, AssetAndInfo)
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

	if name != "" {
		var ok bool
		for _, p := range strings.Split(filepath.ToSlash(name), "/") {
			if node, ok = node[p]; !ok {
				return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
			}
		}
	}

	if len(node) == 0 {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}

	rv := make([]string, 0, len(node))
	for name := range node {
		rv = append(rv, name)
	}

	return rv, nil
}

type bintree map[string]bintree

var _bintree = bintree{
	"1_create_schema_bitburst.down.sql": bintree{},
	"1_create_schema_bitburst.up.sql":   bintree{},
	"2_create_table_objects.down.sql":   bintree{},
	"2_create_table_objects.up.sql":     bintree{},
	"3_create_index_last_seen.down.sql": bintree{},
	"3_create_index_last_seen.up.sql":   bintree{},
}
