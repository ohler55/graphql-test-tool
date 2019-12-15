# GraphQL-Test-Tool

GraphQL Test Tool for running tests cases against a GraphQL server.

## What It Is

The GraphQL Test Tool (gtt) helps test GraphQL servers. Tests are
referred to as use cases and each use case in defined in a JSON
file. Some of the features of gtt are:

 - Both GET and POST HTTP requests can be made.
 - POST content can be either JSON (application/json) or GraphQL (application/graphql).
 - Variables and operation name can be specified in the URL or in JSON content,
 - Values can be remembered and reused in subsequent steps.
 - Various display options.
 - Can be run as an application or the gtt package can be used in unit tests.

## Usage

The GraphQL Test Tool can be run as an application:

```
go run main.go -s http://localhost:6464 -i 2 -v ../examples/top.json
```
The gtt package can be use for unit testing as well. Create the use case files and

```
// TBD ser up runner
uc, err := gtt.NewUseCase("myfile.json"')
```

All tests are driven by use case JSON files. The format is described
in [file_format.md](file_format.md). Some example files are in the
`examples` directory and a simple test server can be set up using the
files in the `test` directory.

## Installation

```
go get github.com/ohler55/graphql-test-tool
```

## Releases

See [CHANGELOG.md](CHANGELOG.md)
