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
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
	"solomatov.me/kuneiform-vscode/lang"
)

type stdioRWC struct{}

type lspHandler struct {
	docs map[string]string
}

func (s *stdioRWC) Close() error {
	return nil
}

func (s *stdioRWC) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (s *stdioRWC) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (l *lspHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "initialize":
		params := lsp.InitializeParams{}
		json.Unmarshal(*req.Params, &params)
		kind := lsp.TDSKFull
		res := lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
				TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
					Kind: &kind,
				},
				DocumentSymbolProvider: true,
			},
		}
		conn.Reply(ctx, req.ID, &res)
	case "textDocument/didOpen":
		params := lsp.DidOpenTextDocumentParams{}
		json.Unmarshal(*req.Params, &params)
		l.docs[string(params.TextDocument.URI)] = params.TextDocument.Text
	case "textDocument/didChange":
		params := lsp.DidChangeTextDocumentParams{}
		json.Unmarshal(*req.Params, &params)
		if len(params.ContentChanges) != 1 {
			panic("Should be exactly one change")
		}
		l.docs[string(params.TextDocument.URI)] = params.ContentChanges[0].Text
	case "textDocument/didClose":
		params := lsp.DidCloseTextDocumentParams{}
		json.Unmarshal(*req.Params, &params)
		delete(l.docs, string(params.TextDocument.URI))
	case "textDocument/documentSymbol":
		params := lsp.DocumentSymbolParams{}
		json.Unmarshal(*req.Params, &params)
		text := l.docs[string(params.TextDocument.URI)]
		f := lang.ParseFile(text)

		fmt.Fprintf(os.Stderr, "File = %#v", f)
	}

}

func main() {
	fmt.Fprintln(os.Stderr, "Starting")
	ctx := context.Background()
	conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(&stdioRWC{}, jsonrpc2.VSCodeObjectCodec{}), &lspHandler{
		docs: map[string]string{},
	})

	fmt.Fprintln(os.Stderr, "Started")
	defer conn.Close()

	<-conn.DisconnectNotify()
}
