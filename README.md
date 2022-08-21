# xgo &nbsp; [![Go Reference](https://pkg.go.dev/badge/github.com/jlourenc/xgo.svg)](https://pkg.go.dev/github.com/jlourenc/xgo) [![GitHub](https://img.shields.io/github/license/jlourenc/xgo)](https://github.com/jlourenc/xgo/blob/main/LICENSE) [![CI](https://github.com/jlourenc/xgo/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/jlourenc/xgo/actions/workflows/ci.yaml) [![Coverage Status](https://coveralls.io/repos/github/jlourenc/xgo/badge.svg?branch=main&service=github)](https://coveralls.io/github/jlourenc/xgo?branch=main) [![CodeQL](https://github.com/jlourenc/xgo/actions/workflows/codeql.yaml/badge.svg?branch=main)](https://github.com/jlourenc/xgo/actions/workflows/codeql.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/jlourenc/xgo)](https://goreportcard.com/report/github.com/jlourenc/xgo)

<p align="center">
  <img alt="xgo logo" src="xgo.png" height="300" />
  <h3 align="center">Extension packages for Go</h3>
</p>

`xgo` is a collection of Go packages that extend the functionalities of the Go standard library. It does not have any external dependencies besides the Go standard library.

Following [Go Release Policy](https://go.dev/doc/devel/release), each major Go release is supported until there are two newer major releases.
Other major releases may still be compatible, however compatibility is not verified nor guaranteed.

# Installation

`xgo` is compatible with modern Go releases in module mode, with Go installed:

```zsh
go get -u github.com/jlourenc/xgo
```

# Contributing

Contributions are very welcomed, big or small!

If you have any questions, please open an issue or a draft PR.

# Versioning

`xgo` follows [semver v2.0.0](https://semver.org/spec/v2.0.0.html) for tagging releases of the package.

* **major version** is incremented with any incompatible change to the exported Go API surface or behavior of the API.
* **minor version** is incremented with any backwards-compatible changes to functionality.
* **patch version** is incremented with any backards-compatible bug fixes.

# License

This project is licensed under the BSD 3-Clause license. See the [LICENSE](/LICENSE) file for more details.
