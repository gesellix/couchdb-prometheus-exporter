Erlang Binary Term Format for Go
================================

[![Build Status](https://secure.travis-ci.org/okeuday/erlang_go.svg?branch=master)](http://travis-ci.org/okeuday/erlang_go) [![Go Report Card](https://goreportcard.com/badge/github.com/okeuday/erlang_go?maxAge=3600)](https://goreportcard.com/report/github.com/okeuday/erlang_go)

Provides all encoding and decoding for the Erlang Binary Term Format
(as defined at [http://erlang.org/doc/apps/erts/erl_ext_dist.html](http://erlang.org/doc/apps/erts/erl_ext_dist.html))
in a single Go package.

(For `go` command-line use you can use the prefix
 `GOPATH=`pwd` GOBIN=$$GOPATH/bin` to avoid additional shell setup)

Build
-----

    go build erlang

Test
----

    go test erlang

Author
------

Michael Truog (mjtruog at protonmail dot com)

License
-------

MIT License
