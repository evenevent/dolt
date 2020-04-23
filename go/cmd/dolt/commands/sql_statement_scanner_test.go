package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestScanStatements(t *testing.T) {
	type testcase struct {
		input string
		statements []string
		lineNums []int
	}

	// Some of these include malformed input (e.g. strings that aren't properly terminated)
	testcases := []testcase {
		{
			input: `insert into foo values (";;';'");`,
			statements: []string{
				`insert into foo values (";;';'")`,
			},
		},
		{
			input: `select ''';;'; select ";\;"`,
			statements: []string{
				`select ''';;'`,
				`select ";\;"`,
			},
		},
		{
			input: `select ''';;'; select ";\;`,
			statements: []string{
				`select ''';;'`,
				`select ";\;`,
			},
		},
		{
			input: `select ''';;'; select ";\;
;`,
			statements: []string{
				`select ''';;'`,
				`select ";\;
;`,
			},
		},
		{
			input: `select '\\'''; select '";";'; select 1`,
			statements: []string{
				`select '\\'''`,
				`select '";";'`,
				`select 1`,
			},
		},
		{
			input: `select '\\''; select '";";'; select 1`,
			statements: []string{
				`select '\\''; select '";"`,
				`'; select 1`,
			},
		},
		{
			input: `insert into foo values(''); select 1`,
			statements: []string{
				`insert into foo values('')`,
				`select 1`,
			},
		},
		{
			input: `insert into foo values('''); select 1`,
			statements: []string{
				`insert into foo values('''); select 1`,
			},
		},
		{
			input: `insert into foo values(''''); select 1`,
			statements: []string{
				`insert into foo values('''')`,
				`select 1`,
			},
		},
		{
			input: `insert into foo values(""); select 1`,
			statements: []string{
				`insert into foo values("")`,
				`select 1`,
			},
		},
		{
			input: `insert into foo values("""); select 1`,
			statements: []string{
				`insert into foo values("""); select 1`,
			},
		},
		{
			input: `insert into foo values(""""); select 1`,
			statements: []string{
				`insert into foo values("""")`,
				`select 1`,
			},
		},
		{
			input: `select '\''; select "hell\"o"`,
			statements: []string{
				`select '\''`,
				`select "hell\"o"`,
			},
		},
		{
			input: `select * from foo; select baz from foo;
select
a from b; select 1`,
			statements: []string{
				"select * from foo",
				"select baz from foo",
				"select\na from b",
				"select 1",
			},
			lineNums: []int {
				1, 1, 2, 3,
			},
		},
		{
			input: "create table dumb (`hell\\`o;` int primary key);",
			statements: []string{
				"create table dumb (`hell\\`o;` int primary key)",
			},
		},
		{
			input: "create table dumb (`hell``o;` int primary key); select \n" +
					"baz from foo;\n" +
					"\n" +
					"select\n" +
					"a from b; select 1\n\n",
			statements: []string{
				"create table dumb (`hell``o;` int primary key)",
				"select \nbaz from foo",
				"select\na from b",
				"select 1",
			},
			lineNums: []int {
				1, 1, 4, 5,
			},
		},
		{
			input: `insert into foo values ('a', "b;", 'c;;""
'); update foo set baz = bar,
qux = '"hello"""' where xyzzy = ";;';'";

  
create table foo (a int not null default ';',
primary key (a));`,
			statements: []string{
				`insert into foo values ('a', "b;", 'c;;""
')`,
				`update foo set baz = bar,
qux = '"hello"""' where xyzzy = ";;';'"`,
				`create table foo (a int not null default ';',
primary key (a))`,
			},
			lineNums: []int {
				1, 2, 6,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			scanner := NewSqlStatementScanner(reader)
			var i int
			for scanner.Scan() {
				require.True(t, i < len(tt.statements))
				assert.Equal(t, tt.statements[i], strings.TrimSpace(scanner.Text()))
				if tt.lineNums != nil {
					assert.Equal(t, tt.lineNums[i], scanner.statementStartLine)
				} else {
					assert.Equal(t, 1, scanner.statementStartLine)
				}
				i++
			}

			require.NoError(t, scanner.Err())
		})
	}
}