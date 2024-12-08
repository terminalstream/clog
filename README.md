[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/terminalstream/clog)
[![Build Status](https://github.com/terminalstream/clog/actions/workflows/ci.yaml/badge.svg)](https://github.com/terminalstream/strum/actions/workflows/ci.yaml?query=branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/terminalstream/clog?style=flat-square)](https://goreportcard.com/report/github.com/terminalstream/clog)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/terminalstream/clog/main/LICENSE)

# [C]ontextual [LOG]ger

`clog` is a very simple wrapper around [`zap.Logger`](https://github.com/uber-go/zap) that
enables [contextual logging](https://github.com/kubernetes/enhancements/blob/master/keps/sig-instrumentation/3077-contextual-logging/README.md).

## Example

```go
package main

import "github.com/terminalstream/clog"

func main() {

}
```
