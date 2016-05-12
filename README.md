gomol-loggly
============

[![GoDoc](https://godoc.org/github.com/aphistic/gomol-loggly?status.svg)](https://godoc.org/github.com/aphistic/gomol-loggly)
[![Build Status](https://img.shields.io/travis/aphistic/gomol-loggly.svg)](https://travis-ci.org/aphistic/gomol-loggly)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/gomol-loggly.svg)](http://codecov.io/github/aphistic/gomol-loggly?branch=master)

gomol-loggly is a logger for [gomol](https://github.com/aphistic/gomol) to support sending logs to [loggly](https://www.loggly.com/).

Installation
============

The recommended way to install is via http://gopkg.in

    go get gopkg.in/aphistic/gomol-loggly.v0
    ...
    import "gopkg.in/aphistic/gomol-loggly.v0"

gomol-loggly can also be installed the standard way as well

    go get github.com/aphistic/gomol-loggly
    ...
    import "github.com/aphistic/gomol-loggly"

Examples
========

For brevity a lot of error checking has been omitted, be sure you do your checks!

This is a super basic example of adding a loggly logger to gomol and then logging a few messages:

```go
package main

import (
	"github.com/aphistic/gomol"
	gl "github.com/aphistic/gomol-loggly"
)

func main() {
	// Add a Loggly logger
	logglyCfg := gl.NewLogglyLoggerConfig()
	logglyCfg.Token = "1234"
	logglyLogger, _ := gl.NewLogglyLogger(logglyCfg)
	gomol.AddLogger(logglyLogger)

	// Set some global attrs that will be added to all
	// messages automatically
	gomol.SetAttr("facility", "gomol.example")
	gomol.SetAttr("another_attr", 1234)

	// Initialize the loggers
	gomol.InitLoggers()
	defer gomol.ShutdownLoggers()

	// Log some debug messages with message-level attrs
	// that will be sent only with that message
	for idx := 1; idx <= 10; idx++ {
		gomol.Dbgm(map[string]interface{}{
			"msg_attr1": 4321,
		}, "Test message %v", idx)
	}
}

```
