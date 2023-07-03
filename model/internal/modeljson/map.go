// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package modeljson

import (
	"go.elastic.co/fastjson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Map struct {
	Struct   *structpb.Struct
	Sanitize bool
}

func (m *Map) MarshalFastJSON(w *fastjson.Writer) error {
	var firstErr error
	w.RawByte('{')
	{
		first := true
		for k, v := range m.Struct.GetFields() {
			if first {
				first = false
			} else {
				w.RawByte(',')
			}
			w.String(k)
			w.RawByte(':')
			if err := fastjson.Marshal(w, v.AsInterface()); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	w.RawByte('}')
	return firstErr
}
