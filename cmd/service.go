/*
Copyright 2023 Rick.

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

package cmd

import (
	"fmt"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	service "github.com/linuxsuren/go-service"
	"github.com/spf13/cobra"
)

type serviceOption struct {
	mode string
}

func newServiceCmd() (c *cobra.Command) {
	opt := &serviceOption{}
	c = &cobra.Command{
		Use:  "service",
		RunE: opt.runE,
		Args: cobra.MinimumNArgs(1),
	}
	flags := c.Flags()
	flags.StringVarP(&opt.mode, "mode", "", "", "The service mode")
	return
}

func (o *serviceOption) runE(c *cobra.Command, args []string) (err error) {
	var svc service.Service
	if svc, err = service.GetAvailableService(service.ServiceMode(o.mode),
		service.ContainerOption{},
		service.CommonService{
			ID:          "env-hub",
			Name:        "env-hub",
			Description: envHubDescription,
			Command:     "env-hub",
			Args:        []string{"server"},
			Execer:      fakeruntime.NewDefaultExecer(),
		}); err != nil {
		return
	}

	var output string
	action := args[0]
	switch action {
	case "start":
		output, err = svc.Start()
	case "stop":
		output, err = svc.Stop()
	case "install":
		output, err = svc.Install()
	case "uninstall":
		output, err = svc.Uninstall()
	case "restart":
		output, err = svc.Restart()
	default:
		err = fmt.Errorf("unknown action %q", action)
	}

	if output != "" {
		c.Println(output)
	}
	return
}
