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

func ParseFile(text string) *FileRoot {
	ctx := newParseCtx(text)
	parseFile(&ctx)
	return ctx.build().(*FileRoot)
}

func parseFile(ctx *parseContext) {
	parseDbDirective(ctx)

	for parseExtDirective(ctx) {
	}

	for parseDecl(ctx) {
	}
}

func parseDbDirective(ctx *parseContext) {
	if ctx.tokKind() != T_DATABASE {
		return
	}
	m := ctx.mark()
	ctx.advance()
	if ctx.tokKind() == T_ID {
		ctx.advance()
	}
	if ctx.tokKind() == T_SEMICOLON {
		ctx.advance()
	}
	m.done(func(ns []AstNode) AstNode { return NewDbDirective(ns) })
}

func parseExtDirective(ctx *parseContext) bool {
	if ctx.tokKind() != T_USE {
		return false
	}
	m := ctx.mark()
	ctx.advance()
	if ctx.tokKind() == T_ID {
		ctx.advance()
	}
	if ctx.tokKind() == T_SEMICOLON {
		ctx.advance()
	}
	m.done(func(ns []AstNode) AstNode { return NewExtDirective(ns) })
	return true
}

func parseDecl(ctx *parseContext) bool {
	if ctx.tokKind() == T_TABLE {
		parseTable(ctx)
		return true
	}

	if ctx.tokKind() == T_ACTION {
		parseAction(ctx)
		return true
	}

	return false

}

func parseTable(ctx *parseContext) bool {
	if ctx.tokKind() != T_TABLE {
		return false
	}

	m := ctx.mark()
	m.ctx.advance()

	if ctx.tokKind() == T_ID {
		ctx.advance()
	}

	if ctx.tokKind() == T_LBRACE {
		ctx.advance()
	}

	if ctx.tokKind() == T_RBRACE {
		ctx.advance()
	}

	m.done(func(ns []AstNode) AstNode { return NewTableDecl(ns) })

	return true
}

func parseAction(ctx *parseContext) bool {
	if ctx.tokKind() != T_ACTION {
		return false
	}

	m := ctx.mark()
	m.ctx.advance()

	if ctx.tokKind() == T_ID {
		ctx.advance()
	}

	if ctx.tokKind() == T_LPAREN {
		ctx.advance()
	}

	for parseParam(ctx) {
		if ctx.tokKind() == T_COMMA {
			ctx.advance()
		}
	}

	if ctx.tokKind() == T_RPAREN {
		ctx.advance()
	}

	if ctx.tokKind() == T_LBRACE {
		ctx.advance()
	}

	for parseStmt(ctx) {
		if ctx.tokKind() == T_SEMICOLON {
			ctx.advance()
		}
	}

	if ctx.tokKind() == T_RBRACE {
		ctx.advance()
	}

	m.done(func(ns []AstNode) AstNode { return NewActionDecl(ns) })

	return true
}

func parseParam(ctx *parseContext) bool {
	if ctx.tokKind() != T_DOLLAR {
		return false
	}

	m := ctx.mark()
	ctx.advance()

	if ctx.tokKind() == T_ID {
		ctx.advance()
	}

	m.done(func(ns []AstNode) AstNode { return NewParamDecl(ns) })

	return true
}

func parseStmt(ctx *parseContext) bool {
	if ctx.tokKind() == T_DOLLAR {
		return parseAssignStmt(ctx)
	}

	return false
}

func parseAssignStmt(ctx *parseContext) bool {
	if ctx.tokKind() != T_DOLLAR {
		return false
	}
	m := ctx.mark()
	ctx.advance()

	if ctx.tokKind() == T_ID {
		ctx.advance()
	}

	if ctx.tokKind() == T_ASSIGN {
		ctx.advance()
		parseExpr(ctx)
	}

	m.done(func(ns []AstNode) AstNode { return NewAssignStmt(ns) })

	return true

}

func parseExpr(ctx *parseContext) bool {
	return parseTermExpr(ctx)
}

func parseTermExpr(ctx *parseContext) bool {
	m := ctx.mark()

	if !parseFactorExpr(ctx) {
		m.drop()
		return false
	}

	for ctx.tokKind() == T_PLUS || ctx.tokKind() == T_MINUS {
		ctx.advance()
		parseFactorExpr(ctx)
		m.done(func(ns []AstNode) AstNode { return NewBinExpr(ns) })
		m = m.precede()
	}

	m.drop()

	return true

}

func parseFactorExpr(ctx *parseContext) bool {
	m := ctx.mark()

	if !parsePrimExpr(ctx) {
		m.drop()
		return false
	}

	for ctx.tokKind() == T_STAR || ctx.tokKind() == T_DIV || ctx.tokKind() == T_MOD {
		ctx.advance()
		parsePrimExpr(ctx)
		m.done(func(ns []AstNode) AstNode { return NewBinExpr(ns) })
		m = m.precede()
	}

	m.drop()

	return true

}

func parsePrimExpr(ctx *parseContext) bool {
	if ctx.tokKind() == T_NUM {
		m := ctx.mark()
		ctx.advance()
		m.done(func(ns []AstNode) AstNode { return NewIntLitExpr(ns) })
		return true
	}

	if ctx.tokKind() == T_DOLLAR {
		m := ctx.mark()
		ctx.advance()
		if ctx.tokKind() == T_ID {
			ctx.advance()
		}
		m.done(func(ns []AstNode) AstNode { return NewVarExpr(ns) })
		return true
	}

	return false
}
