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

import "slices"

type marker struct {
	ctx      *parseContext
	startPos int
	endPos   int
	start    int
	end      int
	dropped  bool
	factory  func([]AstNode) AstNode
}

type parseContext struct {
	tokens  []Token
	markers []*marker
	pos     int
}

func (m *marker) drop() {
	m.dropped = true
}

func (m *marker) precede() *marker {
	if m.dropped {
		panic("Shouldn't be dropped")
	}
	if m.factory == nil {
		panic("Should be completed")
	}

	start := m.start

	res := marker{
		ctx:      m.ctx,
		startPos: m.startPos,
		endPos:   -1,

		start: start,
		end:   -1,

		dropped: false,
		factory: nil,
	}

	for i := len(m.ctx.markers) - 1; i >= start; i-- {
		e := m.ctx.markers[i]
		if i == e.start {
			e.start += 1
		}
		if i == e.end {
			e.end += 1
		}
	}
	m.ctx.markers = slices.Insert(m.ctx.markers, start, &res)
	return &res
}

func (m *marker) rollback() {
	if m.dropped {
		panic("Marker is dropped")
	}
	if m.factory != nil {
		panic("Marker is already done")
	}
	m.ctx.pos = m.startPos
	m.ctx.markers = m.ctx.markers[0:m.start]
}

func (m *marker) done(f func([]AstNode) AstNode) {
	if m.dropped {
		panic("Marker is dropped")
	}
	m.end = len(m.ctx.markers)
	m.endPos = m.ctx.pos - 1
	m.factory = f
	m.ctx.markers = append(m.ctx.markers, m)
}

func (pc *parseContext) skipWs() {
	for {
		k := pc.tokKind()
		if k != T_WS && k != T_COMMENT {
			break
		}
		pc.pos++
	}
}

func (pc *parseContext) tokKind() TokKind {
	if pc.pos >= len(pc.tokens) {
		return T_NONE
	} else {
		return pc.tokens[pc.pos].kind
	}
}

func (pc *parseContext) advance() {
	if pc.pos >= len(pc.tokens) {
		panic("Can't advance")
	}
	pc.pos++
	pc.skipWs()
}

func (pc *parseContext) mark() *marker {
	marker := marker{
		ctx:      pc,
		start:    len(pc.markers),
		end:      -1,
		startPos: pc.pos,
		endPos:   -1,
		dropped:  false,
		factory:  nil,
	}
	pc.markers = append(pc.markers, &marker)
	return &marker
}

func (pc *parseContext) build() AstNode {
	markers := []*marker{nil}
	children := [][]AstNode{{}}

	addChild := func(n AstNode) {
		children[len(children)-1] = append(children[len(children)-1], n)
	}

	tokPos := 0
	for i, m := range pc.markers {
		if m.dropped {
			continue
		}
		if m.factory == nil {
			panic("Inavlid marker")
		}
		switch i {
		case m.start:
			for tokPos < m.startPos {
				addChild(&TokNode{tok: pc.tokens[tokPos]})
				tokPos++
			}
			children = append(children, []AstNode{})
			markers = append(markers, m)
		case m.end:
			if m != markers[len(markers)-1] {
				panic("Marker mismatch")
			}
			for tokPos <= m.endPos {
				addChild(&TokNode{tok: pc.tokens[tokPos]})
				tokPos++
			}
			n := m.factory(children[len(children)-1])
			children = children[0 : len(children)-1]
			markers = markers[0 : len(markers)-1]
			addChild(n)
		default:
			panic("This is impossible")
		}
	}

	if len(markers) != 1 {
		panic("There should be one marker left")
	}

	return NewFileRoot(children[0])
}

func newParseCtx(text string) parseContext {
	toks := tokenize(text)

	res := parseContext{
		tokens: toks,
		pos:    0,
	}
	res.skipWs()
	return res
}
