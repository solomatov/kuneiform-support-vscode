// Copyright 2024 kuneiform-for-vscode contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParsing(t *testing.T) {
	text :=
		`database abc;
		 use xyz;`

	ctx := newParseCtx(text)

	parseFile(&ctx)

	fr := ctx.build().(*FileRoot)

	dd := fr.DbDirective()
	assert.Equal(t, "abc", dd.Name())

	eds := fr.ExtDirectives()
	assert.Equal(t, 1, len(eds))

	assert.Equal(t, "xyz", eds[0].Name())
}

func TestTableDeclParsing(t *testing.T) {
	text :=
		`database abc;
		 table aaa {};`

	ctx := newParseCtx(text)

	parseFile(&ctx)

	fr := ctx.build().(*FileRoot)

	tds := fr.TableDecls()
	assert.Equal(t, 1, len(tds))

	assert.Equal(t, "aaa", tds[0].Name())
}

func TestActionWithParams(t *testing.T) {
	text :=
		`database abc;
		 action bbb ($a, $b, $c) {}`

	ctx := newParseCtx(text)

	parseFile(&ctx)

	fr := ctx.build().(*FileRoot)

	ads := fr.ActionDecls()
	assert.Equal(t, 1, len(ads))

	ad := ads[0]
	assert.Equal(t, "bbb", ad.Name())

	pds := ad.Params()
	assert.Equal(t, 3, len(pds))

	assert.Equal(t, "$a", pds[0].Name())
	assert.Equal(t, "$b", pds[1].Name())
	assert.Equal(t, "$c", pds[2].Name())
}

func TestActionWithBody(t *testing.T) {
	text :=
		`database abc;
		 action bbb ($a) {$a=3;}`

	ctx := newParseCtx(text)

	parseFile(&ctx)

	fr := ctx.build().(*FileRoot)

	ad := fr.ActionDecls()[0]

	sts := ad.Stmts()
	assert.Equal(t, 1, len(sts))
	assert.Equal(t, "$a=3", sts[0].Text())

	as := sts[0].(*AssignStmt)

	e := as.Expr()
	assert.Equal(t, "3", (*e).Text())
}

func TestVarExpr(t *testing.T) {
	e := buildExpr("$x").(*VarExpr)
	assert.Equal(t, "$x", e.VarName())
}

func TestSimpleTermExpr(t *testing.T) {
	e := buildExpr("1+2").(*BinExpr)
	assert.Equal(t, "1+2", e.Text())
}

func TestCorrectAssocTermExpr(t *testing.T) {
	e := buildExpr("1+2+3").(*BinExpr)

	l := *e.Left()
	assert.Equal(t, "1+2", l.Text())

	r := *e.Right()
	assert.Equal(t, "3", r.Text())
}

func TestOpPriority(t *testing.T) {
	e := buildExpr("1+2*3").(*BinExpr)

	l := *e.Left()
	assert.Equal(t, "1", l.Text())

	r := *e.Right()
	assert.Equal(t, "2*3", r.Text())
}

func buildExpr(text string) Expr {
	ctx := newParseCtx("action a(){$x=" + text + ";}")
	parseFile(&ctx)

	fr := ctx.build().(*FileRoot)
	ad := fr.ActionDecls()[0]
	as := ad.Stmts()[0].(*AssignStmt)
	return *as.Expr()
}
