/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

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

package autostart

import (
	"flag"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/vim25/types"
)

type remove struct {
	*AutostartFlag
}

func init() {
	cli.Register("host.autostart.remove", &remove{})
}

func (cmd *remove) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.AutostartFlag, ctx = newAutostartFlag(ctx)
	cmd.AutostartFlag.Register(ctx, f)
}

func (cmd *remove) Process(ctx context.Context) error {
	if err := cmd.AutostartFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *remove) Usage() string {
	return "VM..."
}

func (cmd *remove) Run(ctx context.Context, f *flag.FlagSet) error {
	var powerInfo = types.AutoStartPowerInfo{
		StartAction:      "none",
		StartDelay:       -1,
		StartOrder:       -1,
		StopAction:       "none",
		StopDelay:        -1,
		WaitForHeartbeat: types.AutoStartWaitHeartbeatSettingSystemDefault,
	}

	return cmd.ReconfigureVMs(f.Args(), powerInfo)
}
