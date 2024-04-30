# Kuneiform support for vscode

Jump starting kuneiform support for vscode. This is a work in progress and isn't usable.

## How to setup
* Install node
* Install yarn
* Install go
* 'yarn install' in the root dir
* 'yarn watch' to continously update ext
* 'go build -o server ./server' to build the server

## Useful resources

### Language Services
* Explaining Rust-analyzer video series: https://www.youtube.com/playlist?list=PLhb66M_x9UmrqXhQuIpWC5VgTdrGxMx3y - Detailed explaination of how highly successful language service works internally.
* Blog of the previous videos series author: https://matklad.github.io/

### TMate grammars
* VS Code uses tmate grammars. The best resource is the tmate docs: 
  * Grammar: https://macromates.com/manual/en/language_grammars 
  * Regexps: https://macromates.com/manual/en/regular_expressions
