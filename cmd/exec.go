// Copyright 2022 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package cmd

import (
	"github.com/mendersoftware/go-lib-micro/ws"
	"github.com/mendersoftware/go-lib-micro/ws/exec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack"

	"github.com/mendersoftware/mender-cli/client/deviceconnect"
)

var execCmd = &cobra.Command{
	Use:   "exec device-id \"ls -al | wc -l; touch /tmp/count;\"",
	Short: "Exectue commands on a device",
	Long: "A CLI interface for executing a given command" +
		"and saving the output in the persistent storage as specified by config on device connect",
	Args: cobra.MinimumNArgs(2),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewExec(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

type ExecCmd struct {
	server     string
	skipVerify bool
	token      string

	deviceId string
	command  string
}

func NewExec(cmd *cobra.Command, args []string) (*ExecCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("Empty server value. This should never happen")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	return &ExecCmd{
		server:     server,
		skipVerify: skipVerify,
		token:      token,
		deviceId:   args[0],
		command:    args[1],
	}, nil
}

func (c *ExecCmd) Run() error {
	return c.exec()
}

func (c *ExecCmd) exec() error {

	command := exec.Command{
		CommandLine: c.command,
		Cwd:         "/",
	}
	body, err := msgpack.Marshal(&command)

	client := deviceconnect.NewClient(c.server, c.token, c.skipVerify)

	err = client.Connect(c.deviceId, c.token)
	if err != nil {
		return err
	}

	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:   ws.ProtoTypeExec,
			MsgType: exec.MessageTypeExec,
		},
		Body: body,
	}
	if err = client.WriteMessage(m); err != nil {
		return err
	}
	return nil
}
