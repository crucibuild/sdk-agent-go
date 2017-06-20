/* Copyright (C) 2016 Christophe Camel, Jonathan Pigrée
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agentimpl

import (
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/go-ozzo/ozzo-log"
)

type Logger struct {
	log *log.Logger
}

func NewLogger(a agentiface.Agent) *Logger {
	l := log.NewLogger()
	t1 := log.NewConsoleTarget()
	l.Targets = append(l.Targets, t1)
	l.CallStackDepth = 0

	l.Category = a.Name()
	l.Open()

	logger := &Logger{
		log: l,
	}

	return logger
}

func (l *Logger) Emergency(format string, a ...interface{}) {
	l.log.Emergency(format, a...)
}

func (l *Logger) Alert(format string, a ...interface{}) {
	l.log.Alert(format, a...)
}

func (l *Logger) Critical(format string, a ...interface{}) {
	l.log.Critical(format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.log.Error(format, a...)
}

func (l *Logger) Warning(format string, a ...interface{}) {
	l.log.Warning(format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.log.Info(format, a...)
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.log.Debug(format, a...)
}

func (l *Logger) Close() error {
	l.log.Close()
	return nil
}
