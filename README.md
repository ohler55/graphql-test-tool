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
uc, err := gtt.NewUseCase("myfile.json"')
if err != nil {
    return err
}
runner := gtt.Runner{
    Server: "http://localhost:6464",
    Base:  "/graphql",
    UseCases: []*gtt.UseCase{uc},
}
if err = r.Run(); err != nil {
    return err
}
```

All tests are driven by use case JSON files. The format is described
in [file_format.md](file_format.md). Some example files are in the
`examples` directory and a simple test server can be set up using the
files in the `test` directory.

## Installation

```
go get github.com/ohler55/graphql-test-tool
```

## GoDocs

Documentation is at [https://ohler55.github.io/graphql-test-tool](https://ohler55.github.io/graphql-test-tool).

## Releases

See [CHANGELOG.md](CHANGELOG.md)
