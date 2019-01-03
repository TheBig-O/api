/*
 *   Vikunja is a todo-list application to facilitate your life.
 *   Copyright 2018 Vikunja and contributors. All rights reserved.
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package models

import (
	"code.vikunja.io/web"
	"reflect"
	"runtime"
	"testing"
)

func TestNamespaceUser_Create(t *testing.T) {
	type fields struct {
		ID          int64
		UserID      int64
		NamespaceID int64
		Right       UserRight
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	type args struct {
		a web.Auth
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errType func(err error) bool
	}{
		{
			name: "NamespaceUsers Create normally",
			fields: fields{
				UserID:      1,
				NamespaceID: 2,
			},
		},
		{
			name: "NamespaceUsers Create for duplicate",
			fields: fields{
				UserID:      1,
				NamespaceID: 2,
			},
			wantErr: true,
			errType: IsErrUserAlreadyHasNamespaceAccess,
		},
		{
			name: "NamespaceUsers Create with invalid right",
			fields: fields{
				UserID:      1,
				NamespaceID: 2,
				Right:       500,
			},
			wantErr: true,
			errType: IsErrInvalidUserRight,
		},
		{
			name: "NamespaceUsers Create with inexisting list",
			fields: fields{
				UserID:      1,
				NamespaceID: 2000,
			},
			wantErr: true,
			errType: IsErrNamespaceDoesNotExist,
		},
		{
			name: "NamespaceUsers Create with inexisting user",
			fields: fields{
				UserID:      500,
				NamespaceID: 2,
			},
			wantErr: true,
			errType: IsErrUserDoesNotExist,
		},
		{
			name: "NamespaceUsers Create with the owner as shared user",
			fields: fields{
				UserID:      1,
				NamespaceID: 1,
			},
			wantErr: true,
			errType: IsErrUserAlreadyHasNamespaceAccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &NamespaceUser{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				NamespaceID: tt.fields.NamespaceID,
				Right:       tt.fields.Right,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			err := un.Create(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceUser.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("NamespaceUser.Create() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
		})
	}
}

func TestNamespaceUser_ReadAll(t *testing.T) {
	type fields struct {
		ID          int64
		UserID      int64
		NamespaceID int64
		Right       UserRight
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	type args struct {
		search string
		a      web.Auth
		page   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
		errType func(err error) bool
	}{
		{
			name: "Test readall normal",
			fields: fields{
				NamespaceID: 3,
			},
			args: args{
				a: &User{ID: 3},
			},
			want: []*UserWithRight{
				{
					User: User{
						ID:       1,
						Username: "user1",
						Password: "1234",
						Email:    "user1@example.com",
					},
					Right: UserRightRead,
				},
				{
					User: User{
						ID:       2,
						Username: "user2",
						Password: "1234",
						Email:    "user2@example.com",
					},
					Right: UserRightRead,
				},
			},
		},
		{
			name: "Test ReadAll by a user who does not have access to the list",
			fields: fields{
				NamespaceID: 3,
			},
			args: args{
				a: &User{ID: 4},
			},
			wantErr: true,
			errType: IsErrNeedToHaveNamespaceReadAccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un := &NamespaceUser{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				NamespaceID: tt.fields.NamespaceID,
				Right:       tt.fields.Right,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			got, err := un.ReadAll(tt.args.search, tt.args.a, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceUser.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("NamespaceUser.ReadAll() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceUser.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceUser_Update(t *testing.T) {
	type fields struct {
		ID          int64
		UserID      int64
		NamespaceID int64
		Right       UserRight
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errType func(err error) bool
	}{
		{
			name: "Test Update Normally",
			fields: fields{
				NamespaceID: 3,
				UserID:      1,
				Right:       UserRightAdmin,
			},
		},
		{
			name: "Test Update to write",
			fields: fields{
				NamespaceID: 3,
				UserID:      1,
				Right:       UserRightWrite,
			},
		},
		{
			name: "Test Update to Read",
			fields: fields{
				NamespaceID: 3,
				UserID:      1,
				Right:       UserRightRead,
			},
		},
		{
			name: "Test Update with invalid right",
			fields: fields{
				NamespaceID: 3,
				UserID:      1,
				Right:       500,
			},
			wantErr: true,
			errType: IsErrInvalidUserRight,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nu := &NamespaceUser{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				NamespaceID: tt.fields.NamespaceID,
				Right:       tt.fields.Right,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			err := nu.Update()
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceUser.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("NamespaceUser.Update() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
		})
	}
}

func TestNamespaceUser_Delete(t *testing.T) {
	type fields struct {
		ID          int64
		UserID      int64
		NamespaceID int64
		Right       UserRight
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errType func(err error) bool
	}{
		{
			name: "Try deleting some unexistant user",
			fields: fields{
				UserID:      1000,
				NamespaceID: 2,
			},
			wantErr: true,
			errType: IsErrUserDoesNotExist,
		},
		{
			name: "Try deleting a user which does not has access but exists",
			fields: fields{
				UserID:      1,
				NamespaceID: 4,
			},
			wantErr: true,
			errType: IsErrUserDoesNotHaveAccessToNamespace,
		},
		{
			name: "Try deleting normally",
			fields: fields{
				UserID:      1,
				NamespaceID: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nu := &NamespaceUser{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				NamespaceID: tt.fields.NamespaceID,
				Right:       tt.fields.Right,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			err := nu.Delete()
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceUser.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("NamespaceUser.Delete() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
		})
	}
}