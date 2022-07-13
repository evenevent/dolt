// Copyright 2022 Dolthub, Inc.
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

package docscmds

import (
	"context"

	"github.com/dolthub/dolt/go/cmd/dolt/cli"
	"github.com/dolthub/dolt/go/libraries/doltcore/env"
	"github.com/dolthub/dolt/go/libraries/utils/argparser"
)

var writeDocs = cli.CommandDocumentationContent{
	ShortDesc: "",
	LongDesc:  ``,
	Synopsis:  []string{},
}

type WriteCmd struct{}

// Name implements cli.Command.
func (cmd WriteCmd) Name() string {
	return "write"
}

// Description implements cli.Command.
func (cmd WriteCmd) Description() string {
	return writeDocs.ShortDesc
}

// RequiresRepo implements cli.Command.
func (cmd WriteCmd) RequiresRepo() bool {
	return true
}

// Docs implements cli.Command.
func (cmd WriteCmd) Docs() *cli.CommandDocumentation {
	ap := cmd.ArgParser()
	return cli.NewCommandDocumentation(diffDocs, ap)
}

// ArgParser implements cli.Command.
func (cmd WriteCmd) ArgParser() *argparser.ArgParser {
	ap := argparser.NewArgParser()
	return ap
}

// Exec implements cli.Command.
func (cmd WriteCmd) Exec(ctx context.Context, commandStr string, args []string, dEnv *env.DoltEnv) int {
	panic("dolt docs write")
}
