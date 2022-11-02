// Copyright 2019 Dolthub, Inc.
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

package dtestutils

import (
	"github.com/dolthub/dolt/go/libraries/doltcore/row"
	"github.com/dolthub/dolt/go/libraries/doltcore/schema"
	"github.com/dolthub/dolt/go/store/types"
	"github.com/google/uuid"
)

var UUIDS = []uuid.UUID{
	uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000")),
	uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000001")),
	uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000002"))}
var Names = []string{"Bill Billerson", "John Johnson", "Rob Robertson"}
var Ages = []uint64{32, 25, 21}
var Titles = []string{"Senior Dufus", "Dufus", ""}
var MaritalStatus = []bool{true, false, false}

const (
	IdTag uint64 = iota
	NameTag
	AgeTag
	IsMarriedTag
	TitleTag
	NextTag // leave last
)

const (
	IndexName = "idx_name"
)

var typedColColl = schema.NewColCollection(
	schema.NewColumn("id", IdTag, types.StringKind, true, schema.NotNullConstraint{}),
	schema.NewColumn("name", NameTag, types.StringKind, false, schema.NotNullConstraint{}),
	schema.NewColumn("age", AgeTag, types.UintKind, false, schema.NotNullConstraint{}),
	schema.NewColumn("is_married", IsMarriedTag, types.IntKind, false, schema.NotNullConstraint{}),
	schema.NewColumn("title", TitleTag, types.StringKind, false),
)

// Schema returns the schema for the `people` test table.
func Schema() (schema.Schema, error) {
	sch := schema.MustSchemaFromCols(typedColColl)

	_, err := sch.Indexes().AddIndexByColTags(IndexName, []uint64{NameTag}, schema.IndexProperties{IsUnique: false, Comment: ""})
	if err != nil {
		return nil, err
	}

	_, err = sch.Checks().AddCheck("test-check", "age < 123", true)
	if err != nil {
		return nil, err
	}

	return sch, err
}

// RowsAndSchema returns the schema and rows for the `people` test table.
func RowsAndSchema() ([]row.Row, schema.Schema, error) {
	sch, err := Schema()
	if err != nil {
		return nil, nil, err
	}

	rows := make([]row.Row, len(UUIDS))
	for i := 0; i < len(UUIDS); i++ {
		married := types.Int(0)
		if MaritalStatus[i] {
			married = types.Int(1)
		}
		taggedVals := row.TaggedValues{
			IdTag:        types.String(UUIDS[i].String()),
			NameTag:      types.String(Names[i]),
			AgeTag:       types.Uint(Ages[i]),
			TitleTag:     types.String(Titles[i]),
			IsMarriedTag: married,
		}

		r, err := row.New(types.Format_Default, sch, taggedVals)

		if err != nil {
			panic(err)
		}

		rows[i] = r
	}

	return rows, sch, err
}

// MustMap contructs a types.Tuple for a slice of types.Values.
func MustTuple(vals ...types.Value) types.Tuple {
	tup, err := types.NewTuple(types.Format_Default, vals...)
	if err != nil {
		panic(err)
	}
	return tup
}
