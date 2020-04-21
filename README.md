# Debugger

[![Test](https://github.com/argots/debugger/workflows/Test/badge.svg)](https://github.com/argots/debugger/actions?query=workflow%3ATest)
[![Lint](https://github.com/argots/debugger/workflows/Lint/badge.svg)](https://github.com/argots/debugger/actions?query=workflow%3ALint)
[![Go Report Card](https://goreportcard.com/badge/github.com/argots/debugger)](https://goreportcard.com/report/github.com/argots/debugger)


Debugger helps build debug servers that standard debuggers can connect to.

An implementation of a new language or interpreter or even a DSL can
benefit from debugger support.  This pacakge provides an interface
that new languages can implement to get debugging support via standard
debugger clients (such as Chrome Dev Tools or VS Code).

This is a work in progress:

- [Chrome Dev
Tools](https://github.com/ChromeDevTools/devtools-protocol/) is the
first debug client
- [Debug Adapter
Protocol
(VSCode)](https://microsoft.github.io/debug-adapter-protocol/) is
planned for.

## Chrome Dev Tools

The [fake
server](https://github.com/argots/debugger/tree/master/cmd/fake)
illustrates an example debug server:

```bash
$> go run ./cmd/fake
```

Once that is started, open
[chrome://inspect/#devices](chrome://inspect/#devices), click on "Open
dedicated DevTools for Node" and then type in "localhost:8222".

Now Chrome should be connected to the fake debug server.
