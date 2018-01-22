### Gossock

[![Build Status](https://travis-ci.org/yyyar/gossock.svg?branch=master)](https://travis-ci.org/yyyar/gossock) 
[![Godoc](https://godoc.org/github.com/yyyar/gossock?status.svg)](http://godoc.org/github.com/yyyar/gossock)

Gossock is a small lib for developing lightweight protocols on top of TCP & TLS.

It is an implementation of [Nossock](https://githun.com/yyyar/nossock) in Go.

* **Fast:** serializes messages to JSON, and raw []byte messages
* **Handy interface:** use your structs in simple callback-style handlers
* **TCP and TLS**: easy configurable

#### Installation
```bash
$ go get github.com/yyyar/gossock
```

#### Usage
TODO

#### Tests
```bash
$ go test -v
```

#### Restrictions
Due to message structure, message body size is limited to `4096 Mb`. No idea why
you'll want to send such a big message, in any case it worth to split it to
lots of smaller parts.

#### Author
* [Yaroslav Pogrebnyak](https://github.com/yyyar/)

#### License
MIT

