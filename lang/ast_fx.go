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

type AstNode interface {
	Children() []AstNode
	TextLen() int
	Text() string
}

type TokNode struct {
	tok Token
}

type compNode struct {
	len   int
	nodes []AstNode
}

func (t *TokNode) Children() []AstNode {
	return []AstNode{}
}

func (t *TokNode) TextLen() int {
	return len(t.tok.text)
}

func (t *TokNode) Text() string {
	return t.tok.text
}

func newTok(t Token) *TokNode {
	return &TokNode{
		tok: t,
	}
}

func (c *compNode) Children() []AstNode {
	return c.nodes
}

func (c *compNode) TextLen() int {
	return c.len
}

func (c *compNode) Text() string {
	res := ""
	for _, c := range c.nodes {
		res += c.Text()
	}
	return res
}

func newComp(ns []AstNode) *compNode {
	tl := 0
	for _, n := range ns {
		tl += n.TextLen()
	}

	return &compNode{
		nodes: ns,
		len:   tl,
	}
}

func findTok(ns []AstNode, tk TokKind) *TokNode {
	for _, n := range ns {
		t, ok := n.(*TokNode)
		if ok && t.tok.kind == tk {
			return t
		}
	}
	return nil
}

func idText(ns []AstNode) string {
	id := findTok(ns, T_ID)
	if id == nil {
		return ""
	} else {
		return id.Text()
	}
}

func typedChildren[A AstNode](ns []AstNode) []A {
	res := []A{}
	for _, c := range ns {
		i, ok := c.(A)
		if ok {
			res = append(res, i)
		}
	}
	return res
}
