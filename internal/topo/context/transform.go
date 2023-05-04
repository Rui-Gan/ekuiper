// Copyright 2021-2023 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package context

import (
	"fmt"
	"github.com/lf-edge/ekuiper/internal/topo/transform"
)

const TransKey = "$$trans"
const SelectKey = "$$select"
const SelectFields = "$$selectFields"

// TransformOutput Lazy transform output to bytes
func (c *DefaultContext) TransformOutput(data interface{}) ([]byte, bool, error) {
	v := c.Value(TransKey)
	f, ok := v.(transform.TransFunc)
	if !ok {
		return nil, false, fmt.Errorf("no transform configured")
	}

	bytes, ok, err := f(data)
	if err != nil || ok || c.Value(SelectFields) == nil {
		return bytes, ok, err
	}

	fields := c.Value(SelectFields).([]string)
	s := c.Value(SelectKey)
	sf, ok := s.(transform.SelectFunc)
	if !ok {
		return bytes, false, fmt.Errorf("no select configured")
	}
	return sf(bytes, fields)
}
