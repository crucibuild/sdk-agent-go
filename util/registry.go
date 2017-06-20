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

package util

import "github.com/cheekybits/genny/generic"

type KEY generic.Type
type VALUE generic.Type
type NAME generic.Type

type NAMERegistry interface {
	// RegisterNAME registers a new NAME in the registry.
	RegisterNAME(item VALUE) (KEY, error)
	// GetNAME retrieve the NAME in the registry given its key.
	GetNAME(key KEY) (VALUE, error)
	// ListNAMEs list all the items (NAMEs) contained in the registry.
	ListNAMEs() []VALUE
	// UnregisterNAME removes the item from the registry.
	UnregisterNAME(key KEY) error
	// ExistNAME check if the given key exists in the registry.
	ExistNAME(key KEY) bool
}
