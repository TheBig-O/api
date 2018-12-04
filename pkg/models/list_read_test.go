//  Vikunja is a todo-list application to facilitate your life.
//  Copyright 2018 Vikunja and contributors. All rights reserved.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestList_ReadAll(t *testing.T) {
	// Create test database
	//assert.NoError(t, PrepareTestDatabase())

	// Get all lists for our namespace
	lists, err := GetListsByNamespaceID(1, &User{})
	assert.NoError(t, err)
	assert.Equal(t, len(lists), 2)

	// Get all lists our user has access to
	u, err := GetUserByID(1)
	assert.NoError(t, err)

	lists2 := List{}
	lists3, err := lists2.ReadAll("", &u, 1)
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(lists3).Kind(), reflect.Slice)
	s := reflect.ValueOf(lists3)
	assert.Equal(t, s.Len(), 1)

	// Try getting lists for a nonexistant user
	_, err = lists2.ReadAll("", &User{ID: 984234}, 1)
	assert.Error(t, err)
	assert.True(t, IsErrUserDoesNotExist(err))
}