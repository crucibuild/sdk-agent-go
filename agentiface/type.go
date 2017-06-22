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

package agentiface

import "reflect"

type Type interface {
	Name() string
	Type() reflect.Type
}

type TypeRegistry interface {
	// TypeRegister registers a new Type in the registry.
	TypeRegister(t Type) (string, error)
	// TypeGetByName retrieve the Type in the registry given its key.
	TypeGetByName(key string) (Type, error)
	// TypeGetByType retrieve the type in the registry given its type.
	TypeGetByType(t reflect.Type) (Type, error)
	// TypeListNames list all the item names (Types) contained in the registry.
	TypeListNames() []string
	// TypeUnregister removes the item from the registry.
	TypeUnregister(key string) error
	// TypeExist check if the given key exists in the registry.
	TypeExist(key string) bool
}
