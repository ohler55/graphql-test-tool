# Building Solid Go GraphQL Applications Quickly

With the experience gained from developing over half a dozen deployed
GraphQL application this article is a tour of a testing methodology
that has been key to developing solid GraphQL applications
quickly. The primary contributors to success are a good GraphQL server
package to build the application and suitable black box testing tool.

## Black Box Testing

It would be hard to argue that application level black box testing is
not the gold standard for testing. Sure, unit testing is useful and
necessary but the final gate before release or deployment really
should be tests that validate behavior as the end users will see it.

GraphQL servers are no exception to the black box testing rule. A test
suite that exercises an GraphQL application using the same HTTP API
that end users will use make it less likely there will be surprises
after deployment. A black box test suite also lends itself to
continuous integration tests on merges during development.

With the proper tooling testing through the HTTP API of a GraphQL
application is often easier than trying to piece together regression
tests that only exercise internal APIs. To be effective the tool
should be able to play the role of a
user. [GraphQL-Test-Tool](https://github.com/ohler55/graphql-test-tool)
is such a tool. GTT for short,
[GraphQL-Test-Tool](https://github.com/ohler55/graphql-test-tool) is a
script driven test tool for GraphQL. I addition to providing
repeatable regression tests and CI GTT scripts are great examples for
end users.

GTT is written in Go but can be used for testing any GraphQL server
that has an HTTP API. The
[example](https://github.com/ohler55/graphql-test-tool/tree/master/example)
demonstrates a Go application test as well as a Ruby server test but
the focus of this article is on the Go test setup.

## Application

The example GraphQL application is taken from the
[GGql](https://github.com/UHN/ggql) reflection example. GGql is the
fastest Go GraphQL server as well as being the easiest to use as shown
by this
[comparison](https://github.com/the-benchmarker/graphql-benchmarks/blob/develop/rates.md).
Instead of going into detail here of how the application was put
together is would be better to read the details in the
[README.md](https://github.com/UHN/ggql/tree/master/examples/reflection/README.md)
file for the example. The only changes to that example were the
addition of a command option to set the port and support for the
`/graphql/schema` URL path to return the schema as SDL.

A [GGql](https://github.com/UHN/ggql) root object is able to return
the formatted schema that it is serving. By adding an HTTP handler to
`/graphql/schema` an HTTP request can be made to respond with the full
schema in SDL format.

``` golang
	http.HandleFunc("/graphql/schema", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		full := strings.EqualFold(q.Get("full"), "true")
		desc := strings.EqualFold(q.Get("desc"), "true")
		sdl := root.SDL(full, desc)
		_, _ = w.Write([]byte(sdl))
	})
```

## Test Setup

With the application ready to go a test setup is next. There are two
choices when setting up the tests. One is to use the full application
for a true black box test and the other is to split the application
into a command portion and a package portion so that the package
portion can be tested in the same process space as the test
code. There are advantages and disadvantages to each.

### True Black Box

A true black box approach runs the application completely separate
from the test code. This had the advantage of being able to test a
server implemented in any language. The disadvantage is that the Go
test coverage tool will not work. Running as a separate process also
means debugging pront statements are a little more difficult to
display.

### Embedded Tests

In order to run the application in the same code space as the test
code it needs to be callable from the test code which means it needs
to be in a package that can be imported. Thats not difficult to
do. Just create a `cmd` directory and put a light weight `main()`
function in an application directory of the `cmd` directory that calls
the package where all the rest of the code resides. Go coverage tools
then work and debug print statements show up as the tests are
running. The downside to an embedded configuration is that the server
has to be written in Go which isn't really much of a downside.

### Implementation

Since the true black box is a bit tricker to set up and since it also
allows for testing servers written in languages other than Go that is
the path taken for this article.

Test scripts are placed in a `gtt` subdirectory for cleanliness only.

While not necessary, the application will be started just once and
then each test will be executed against the running app. The code to
run the application is in `main_test.go` while the test functions for
running the tests are in `song_test.go`.

#### `main_test.go`


#### `song_test.go`



## Test Scripts

Script files are in a separate directory named `gtt`. More details on
the script file are in
[file_format.md](https://github.com/ohler55/graphql-test-tool/blob/master/file_format.md)
which can be referred to while we walk though the scripts.

 - TBD

## Summary

While other packages were involved in managing the data behind the
GraphQL servers, [GGql](https://github.com/UHN/ggql) was key in
building the applications
quickly. [GTT](https://github.com/ohler55/graphql-test-tool) was key
to testing and providing examples to the teams writing the user
portions of the project.
