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
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokKind string

const (
	T_COLON     TokKind = ":"
	T_SEMICOLON TokKind = ";"
	T_LPAREN    TokKind = "("
	T_RPAREN    TokKind = ")"
	T_LBRACE    TokKind = "{"
	T_RBRACE    TokKind = "}"
	T_COMMA     TokKind = ","
	T_DOLLAR    TokKind = "$"
	T_HASH      TokKind = "#"
	T_AT        TokKind = "@"
	T_DOT       TokKind = "."
	T_ASSIGN    TokKind = "="
	T_PLUS      TokKind = "+"
	T_MINUS     TokKind = "-"
	T_STAR      TokKind = "*"
	T_DIV       TokKind = "/"
	T_MOD       TokKind = "%"
	T_TILDE     TokKind = "~"
	T_LOGIC_OR  TokKind = "||"
	T_LSHIFT    TokKind = "<<"
	T_RSHIFT    TokKind = ">>"
	T_AND       TokKind = "&"
	T_OR        TokKind = "|"
	T_EQ        TokKind = "=="
	T_LESS      TokKind = "<"
	T_LESS_EQ   TokKind = "<="
	T_GT        TokKind = ">"
	T_GT_EQ     TokKind = ">="
	T_NOT_EQ    TokKind = "!="
	T_NEQ       TokKind = "<>"

	T_ID      TokKind = "id"
	T_WS      TokKind = "ws"
	T_COMMENT TokKind = "comment"
	T_NUM     TokKind = "num"

	T_DATABASE TokKind = "database"
	T_USE      TokKind = "use"
	T_TABLE    TokKind = "table"
	T_ACTION   TokKind = "action"

	T_ERROR TokKind = "error"
	T_NONE  TokKind = "none"
)

type Token struct {
	start int
	end   int
	text  string
	kind  TokKind
}

func tokenize(text string) []Token {
	res := []Token{}

	cur := 0
	tokStart := 0
	curRune := func() rune {
		r, _ := utf8.DecodeRuneInString(text[cur:])
		return r
	}
	advance := func() {
		r, l := utf8.DecodeRuneInString(text[cur:])
		if r == utf8.RuneError {
			panic("Can't advance")
		}
		cur += l
	}
	tokText := func() string {
		return text[tokStart:cur]
	}
	finish := func(kind TokKind) {
		if cur == tokStart {
			return
		}
		tok := Token{
			start: tokStart,
			end:   cur,
			text:  tokText(),
			kind:  kind,
		}
		tokStart = cur
		res = append(res, tok)
	}

	for {
		r := curRune()
		if r == utf8.RuneError {
			break
		}

		switch r {

		case ':', ';', '(', ')', '{', '}', ',', '.', '$', '#', '@', '+', '-', '*', '%', '~', '&':
			advance()
			finish(TokKind(string(r)))

		case '/':
			advance()
			if curRune() == '/' {
				advance()
				for curRune() != '\r' && curRune() != '\n' && curRune() != utf8.RuneError {
					advance()
				}
				finish(T_COMMENT)
			} else if curRune() == '*' {
				advance()
				for {
					if curRune() == utf8.RuneError {
						finish(T_COMMENT)
						break
					}
					if curRune() == '*' {
						advance()
						if curRune() == '/' {
							advance()
							finish(T_COMMENT)
							break
						} else if curRune() == utf8.RuneError {
							finish(T_COMMENT)
							break
						}
					}
					advance()
				}
			} else {
				finish(T_DIV)
			}

		case '=':
			advance()
			if curRune() == '=' {
				advance()
				finish(T_EQ)
			} else {
				finish(T_ASSIGN)
			}

		case '|':
			advance()
			if curRune() == '|' {
				advance()
				finish(T_LOGIC_OR)
			} else {
				finish(T_OR)
			}
		case '!':
			advance()
			if curRune() == '=' {
				advance()
				finish(T_NOT_EQ)
			}
		case '<':
			advance()
			if curRune() == '<' {
				advance()
				finish(T_LSHIFT)
			} else if curRune() == '=' {
				advance()
				finish(T_LESS_EQ)
			} else if curRune() == '>' {
				advance()
				finish(T_NEQ)
			} else {
				finish(T_LESS)
			}
		case '>':
			advance()
			if curRune() == '=' {
				advance()
				finish(T_LESS_EQ)
			} else if curRune() == '>' {
				advance()
				finish(T_RSHIFT)
			} else {
				advance()
				finish(T_GT)
			}
		default:
			if unicode.IsLetter(r) {
				advance()
				for IsLetter(curRune()) || unicode.IsDigit(curRune()) {
					advance()
				}

				tokText := tokText()
				switch strings.ToLower(tokText) {
				case "database":
					finish(T_DATABASE)
				case "use":
					finish(T_USE)
				case "table":
					finish(T_TABLE)
				case "action":
					finish(T_ACTION)
				default:
					finish(T_ID)
				}
			} else if unicode.IsDigit(r) {
				advance()
				for unicode.IsDigit(curRune()) {
					advance()
				}
				finish(T_NUM)
			} else if unicode.IsSpace(curRune()) {
				advance()
				for unicode.IsSpace(curRune()) {
					advance()
				}
				finish(T_WS)
			} else {
				advance()
				finish(T_ERROR)
			}
		}
	}

	finish(T_ERROR)

	return res
}

func IsLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func IsSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}
