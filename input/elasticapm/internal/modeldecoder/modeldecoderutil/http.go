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

package modeldecoderutil

import (
	"net/http"

	"google.golang.org/protobuf/types/known/structpb"
)

// HTTPHeadersToMap converts h to a map[string]any, suitable for
// use in model.HTTP.{Request,Response}.Headers.
func HTTPHeadersToMap(h http.Header) *structpb.Struct {
	if len(h) == 0 {
		return nil
	}
	m := make(map[string]any, len(h))
	for k, v := range h {
		arr := make([]any, 0, len(v))
		for _, s := range v {
			arr = append(arr, s)
		}
		m[k] = arr
	}
	if m, err := structpb.NewStruct(m); err == nil {
		return m
	}
	return nil
}

// NormalizeHTTPRequestBody recurses through v, replacing any instance of
// a json.Number with float64.
//
// TODO(axw) define a more restrictive schema for context.request.body
// so this is unnecessary. Agents are unlikely to send numbers, but
// seeing as the schema does not prevent it we need this.
func NormalizeHTTPRequestBody(v interface{}) *structpb.Value {
	if v, err := structpb.NewValue(v); err == nil {
		return v
	}
	return nil
}
