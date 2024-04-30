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

type FileRoot struct {
	compNode
}

type DbDirective struct {
	compNode
}

type ExtDirective struct {
	compNode
}

type TableDecl struct {
	compNode
}

type ActionDecl struct {
	compNode
}

type ParamDecl struct {
	compNode
}

type Stmt interface {
	AstNode
	IsStmt()
}

type AssignStmt struct {
	compNode
}

type Expr interface {
	AstNode
	IsExpr()
}

type VarExpr struct {
	compNode
}

type BinExpr struct {
	compNode
}

type IntLitExpr struct {
	compNode
}

func (fr *FileRoot) DbDirective() *DbDirective {
	for _, c := range fr.Children() {
		r, ok := c.(*DbDirective)
		if ok {
			return r
		}
	}
	return nil
}

func (fr *FileRoot) ExtDirectives() []*ExtDirective {
	return typedChildren[*ExtDirective](fr.Children())
}

func (fr *FileRoot) TableDecls() []*TableDecl {
	return typedChildren[*TableDecl](fr.Children())
}

func (fr *FileRoot) ActionDecls() []*ActionDecl {
	return typedChildren[*ActionDecl](fr.Children())
}

func (dd *DbDirective) Name() string {
	return idText(dd.Children())
}

func (ed *ExtDirective) Name() string {
	return idText(ed.Children())
}

func (td *TableDecl) Name() string {
	return idText(td.Children())
}

func (ad *ActionDecl) Name() string {
	return idText(ad.Children())
}

func (ad *ActionDecl) Params() []*ParamDecl {
	return typedChildren[*ParamDecl](ad.Children())
}

func (ad *ActionDecl) Stmts() []Stmt {
	return typedChildren[Stmt](ad.Children())
}

func (pd *ParamDecl) Name() string {
	return "$" + idText(pd.Children())
}

func (as *AssignStmt) IsStmt() {}

func (as *AssignStmt) Expr() *Expr {
	exprs := typedChildren[Expr](as.Children())
	if len(exprs) == 0 {
		return nil
	} else {
		return &exprs[0]
	}
}

func (ve *VarExpr) IsExpr() {}

func (ve *VarExpr) VarName() string {
	return "$" + idText(ve.Children())
}

func (be *BinExpr) IsExpr() {}

func (be *BinExpr) Left() *Expr {
	exprs := typedChildren[Expr](be.Children())
	if len(exprs) > 0 {
		return &exprs[0]
	} else {
		return nil
	}
}

func (be *BinExpr) Right() *Expr {
	exprs := typedChildren[Expr](be.Children())
	if len(exprs) > 1 {
		return &exprs[1]
	} else {
		return nil
	}
}

func (il *IntLitExpr) IsExpr() {}

func NewFileRoot(ns []AstNode) *FileRoot {
	return &FileRoot{
		compNode: *newComp(ns),
	}
}

func NewDbDirective(ns []AstNode) *DbDirective {
	return &DbDirective{
		compNode: *newComp(ns),
	}
}

func NewExtDirective(ns []AstNode) *ExtDirective {
	return &ExtDirective{
		compNode: *newComp(ns),
	}
}

func NewTableDecl(ns []AstNode) *TableDecl {
	return &TableDecl{
		compNode: *newComp(ns),
	}
}

func NewActionDecl(ns []AstNode) *ActionDecl {
	return &ActionDecl{
		compNode: *newComp(ns),
	}
}

func NewParamDecl(ns []AstNode) *ParamDecl {
	return &ParamDecl{
		compNode: *newComp(ns),
	}
}

func NewAssignStmt(ns []AstNode) *AssignStmt {
	return &AssignStmt{
		compNode: *newComp(ns),
	}
}

func NewVarExpr(ns []AstNode) *VarExpr {
	return &VarExpr{
		compNode: *newComp(ns),
	}
}

func NewBinExpr(ns []AstNode) *BinExpr {
	return &BinExpr{
		compNode: *newComp(ns),
	}
}

func NewIntLitExpr(ns []AstNode) *IntLitExpr {
	return &IntLitExpr{
		compNode: *newComp(ns),
	}
}
