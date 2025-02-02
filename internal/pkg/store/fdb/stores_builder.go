// Copyright 2023 EMQ Technologies Co., Ltd.
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

//go:build amd64 && (fdb || full)

package fdb

import "github.com/lf-edge/ekuiper/internal/pkg/store/definition"

func BuildStores(c definition.Config, name string) (definition.StoreBuilder, definition.TsBuilder, error) {
	db, err := NewFdbFromConf(c)
	if err != nil {
		return nil, nil, err
	}
	kvBuilder := NewStoreBuilder(db)
	tsBuilder := NewTsBuilder(db)
	return kvBuilder, tsBuilder, nil
}
