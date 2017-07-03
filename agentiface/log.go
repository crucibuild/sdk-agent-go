// Copyright (C) 2016 Christophe Camel, Jonathan Pigrée
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agentiface

type Logger interface {
	// Emergency logs a message indicating the system is unusable.
	Emergency(format string, a ...interface{})
	// Alert logs a message indicating action must be taken immediately.
	Alert(format string, a ...interface{})
	// Critical logs a message indicating critical conditions.
	Critical(format string, a ...interface{})
	// Error logs a message indicating an error condition.
	// This method takes one or multiple parameters. If a single parameter
	// is provided, it will be treated as the log message. If multiple parameters
	// are provided, they will be passed to fmt.Sprintf() to generate the log message.
	Error(format string, a ...interface{})
	// Warning logs a message indicating a warning condition.
	Warning(format string, a ...interface{})
	// Info logs a message for informational purpose.
	Info(format string, a ...interface{})
	// Debug logs a message for debugging purpose.
	Debug(format string, a ...interface{})
}
