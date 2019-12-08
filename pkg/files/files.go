// Vikunja is a todo-list application to facilitate your life.
// Copyright 2019 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package files

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/web"
	"github.com/c2h5oh/datasize"
	"github.com/spf13/afero"
	"io"
	"strconv"
	"time"
)

// File holds all information about a file
type File struct {
	ID   int64  `xorm:"int(11) autoincr not null unique pk" json:"id"`
	Name string `xorm:"text not null" json:"name"`
	Mime string `xorm:"text null" json:"mime"`
	Size uint64 `xorm:"int(11) not null" json:"size"`

	Created time.Time `xorm:"-" json:"created"`

	CreatedUnix int64 `xorm:"created" json:"-"`
	CreatedByID int64 `xorm:"int(11) not null" json:"-"`

	File afero.File `xorm:"-" json:"-"`
}

// TableName is the table name for the files table
func (File) TableName() string {
	return "files"
}

func (f *File) getFileName() string {
	return config.FilesBasePath.GetString() + "/" + strconv.FormatInt(f.ID, 10)
}

// LoadFileByID returns a file by its ID
func (f *File) LoadFileByID() (err error) {
	f.File, err = afs.Open(f.getFileName())
	return
}

// LoadFileMetaByID loads everything about a file without loading the actual file
func (f *File) LoadFileMetaByID() (err error) {
	exists, err := x.Where("id = ?", f.ID).Get(f)
	if !exists {
		return ErrFileDoesNotExist{FileID: f.ID}
	}
	f.Created = time.Unix(f.CreatedUnix, 0)
	return
}

// Create creates a new file from an FileHeader
func Create(f io.ReadCloser, realname string, realsize uint64, a web.Auth) (file *File, err error) {

	// Get and parse the configured file size
	var maxSize datasize.ByteSize
	err = maxSize.UnmarshalText([]byte(config.FilesMaxSize.GetString()))
	if err != nil {
		return nil, err
	}
	if realsize > maxSize.Bytes() {
		return nil, ErrFileIsTooLarge{Size: realsize}
	}

	// We first insert the file into the db to get it's ID
	file = &File{
		Name:        realname,
		Size:        realsize,
		CreatedByID: a.GetID(),
	}

	_, err = x.Insert(file)
	if err != nil {
		return
	}

	// Save the file to storage with its new ID as path
	err = afs.WriteReader(file.getFileName(), f)
	return
}

// Delete removes a file from the DB and the file system
func (f *File) Delete() (err error) {
	deleted, err := x.Where("id = ?", f.ID).Delete(f)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return ErrFileDoesNotExist{FileID: f.ID}
	}

	err = afs.Remove(f.getFileName())
	return
}