// Copyright 2022 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package exec

const (
	// MessageTypeExec requests a command execution on a device.
	MessageTypeExec = "exec"
	// MessageTypeError is returned on internal or protocol errors. The
	// body MUST contain an Error object.
	MessageTypeError = "error"
)

// The Error struct is passed in the Body of MsgProto in case the message type is ErrorMessage
type Error struct {
	// The error description, as in "Permission denied while opening a file"
	Error *string `msgpack:"err" json:"error"`
	// Type of message that raised the error
	MessageType *string `msgpack:"msgtype,omitempty" json:"message_type,omitempty"`
	// Message id is passed in the MsgProto Properties, and in case it is available and
	// error occurs it is passed for reference in the Body of the error message
	MessageID *string `msgpack:"msgid,omitempty" json:"message_id,omitempty"`
}

type Command struct {
	CommandLine string
	Cwd string
}