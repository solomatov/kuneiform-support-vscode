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

func TestSimpleLexing(t *testing.T) {
	text := "<<+>>"
	toks := tokenize(text)

	assert.Equal(t, []Token{
		{
			start: 0,
			end:   2,
			text:  "<<",
			kind:  T_LSHIFT,
		},
		{
			start: 2,
			end:   3,
			text:  "+",
			kind:  T_PLUS,
		},
		{
			start: 3,
			end:   5,
			text:  ">>",
			kind:  T_RSHIFT,
		},
	}, toks)

}

func TestLessAmbiguity(t *testing.T) {
	text := "<<<"
	toks := tokenize(text)

	assert.Equal(t, []Token{
		{
			start: 0,
			end:   2,
			text:  "<<",
			kind:  T_LSHIFT,
		},
		{
			start: 2,
			end:   3,
			text:  "<",
			kind:  T_LESS,
		},
	}, toks)
}

func TestSpace(t *testing.T) {
	text := "ab \n t"
	toks := tokenize(text)

	assert.Equal(t, []Token{
		{
			start: 0,
			end:   2,
			text:  "ab",
			kind:  T_ID,
		},
		{
			start: 2,
			end:   5,
			text:  " \n ",
			kind:  T_WS,
		},
		{
			start: 5,
			end:   6,
			text:  "t",
			kind:  T_ID,
		},
	}, toks)

}

func TestLineComment(t *testing.T) {
	text := "ab\n//aaaa\nb"
	toks := tokenize(text)

	assert.Equal(t, []Token{
		{
			start: 0,
			end:   2,
			text:  "ab",
			kind:  T_ID,
		},
		{
			start: 2,
			end:   3,
			text:  "\n",
			kind:  T_WS,
		},
		{
			start: 3,
			end:   9,
			text:  "//aaaa",
			kind:  T_COMMENT,
		},
		{
			start: 9,
			end:   10,
			text:  "\n",
			kind:  T_WS,
		},
		{
			start: 10,
			end:   11,
			text:  "b",
			kind:  T_ID,
		},
	}, toks)
}

func TestMultiline(t *testing.T) {
	text := " /*a*/b"
	toks := tokenize(text)

	assert.Equal(t, []Token{
		{
			start: 0,
			end:   1,
			text:  " ",
			kind:  T_WS,
		},
		{
			start: 1,
			end:   6,
			text:  "/*a*/",
			kind:  T_COMMENT,
		},
		{
			start: 6,
			end:   7,
			text:  "b",
			kind:  T_ID,
		},
	}, toks)

}

func TestIncompleteMultiLineComment(t *testing.T) {
	text := "/*aaaa*"
	toks := tokenize(text)

	assert.Equal(t, []Token{
		{
			start: 0,
			end:   7,
			text:  "/*aaaa*",
			kind:  T_COMMENT,
		},
	}, toks)
}
