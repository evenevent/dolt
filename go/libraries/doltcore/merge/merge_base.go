// Copyright 2021 Dolthub, Inc.
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

package merge

import (
	"context"

	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb"

	"github.com/dolthub/dolt/go/store/hash"
)

func MergeBase(ctx context.Context, left, right *doltdb.Commit) (base hash.Hash, err error) {
	optCmt, err := doltdb.GetCommitAncestor(ctx, left, right)
	if err != nil {
		return base, err
	}
	ancestor, ok := optCmt.ToCommit()
	if !ok {
		return base, doltdb.ErrGhostCommitEncountered // NM4 - TEST.  I think getCommitAncestor is going to be an awk one.
	}

	return ancestor.HashOf()
}
