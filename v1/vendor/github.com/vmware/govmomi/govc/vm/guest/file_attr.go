/*
Copyright (c) 2014-2016 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package guest

import (
	"flag"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/vim25/types"
)

type FileAttrFlag struct {
	types.GuestPosixFileAttributes
}

func newFileAttrFlag(ctx context.Context) (*FileAttrFlag, context.Context) {
	return &FileAttrFlag{}, ctx
}

func (flag *FileAttrFlag) Register(ctx context.Context, f *flag.FlagSet) {
	f.Var(flags.NewInt32(&flag.OwnerId), "uid", "User ID")
	f.Var(flags.NewInt32(&flag.GroupId), "gid", "Group ID")
	f.Int64Var(&flag.Permissions, "perm", 0, "File permissions")
}

func (flag *FileAttrFlag) Process(ctx context.Context) error {
	return nil
}

func (flag *FileAttrFlag) Attr() types.BaseGuestFileAttributes {
	return &flag.GuestPosixFileAttributes
}
