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

Since the tests can not be run without first starting the application
it will be necessary to create a `TestMain()` function. After setting
up a flag for executing a Ruby version of the application a call to
the `run()` function is made. That function will do most of the work
and return an error if anything goes wrong with the setup or the tests
fail.

``` golang
func TestMain(m *testing.M) {
	flag.BoolVar(&ruby, "ruby", ruby, "run the ruby server instead of go server")
	flag.Parse()

	if err := run(m); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
```

The `run()` function first finds a free port to run the application on.

``` golang
func run(m *testing.M) (err error) {
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var ln *net.TCPListener
		if ln, err = net.ListenTCP("tcp", addr); err == nil {
			testPort = ln.Addr().(*net.TCPAddr).Port
			ln.Close()
		}
		if err != nil {
			return
		}
	}
```

Then the application is started. Since we want to collect the output
from the application we grab the application stdout before starting.

``` golang
	var cmd *exec.Cmd
	if ruby {
		cmd = exec.Command("ruby", "song.rb", "-p", strconv.Itoa(testPort))
	} else {
		cmd = exec.Command("go", "run", "main.go", "-p", strconv.Itoa(testPort))
	}
	stdout, _ := cmd.StdoutPipe()
	if err = cmd.Start(); err != nil {
		return
	}
```

A go routine running concurrently will have to read the stdout while
the application is running. The buffer size of stdout is limited so if
the application generates too much output the application will hang
trying to write to stdout. With the reading being done concurrently
the `ioutil.ReadAll()` buffer will expand and collect all the
output. After stdout is closed `ioutil.ReadAll()` will return but the
main thread needs to be told it is free to continue and print the
collected output. The `done` chan takes care of that.

``` golang
	var out []byte
	done := make(chan bool)
	go func() {
		out, _ = ioutil.ReadAll(stdout)
		done <- true
	}()
```

The tests really shouldn't be started until the application has
started. We could sleep but then that forces a slowdown of the
test. It's better to continue once the application is up and accepting
requests. A loop with a delay between attempts takes care of that.

``` golang
	for i := 0; i < 25; i++ {
		u := fmt.Sprintf("http://localhost:%d", testPort)
		var r *http.Response
		if r, err = http.Get(u); err == nil {
			r.Body.Close()
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
```

If there were no errors the tests can be run.

``` golang
	if err == nil && 0 != m.Run() {
		err = fmt.Errorf("tests failed")
	}
}
```

After the tests finish it's time to cleanup but killing the
application and printing the application output. Note the wait on the
`done` chan before proceeding to printing.

``` golang
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	stdout.Close()
	<-done
	if testing.Verbose() {
		fmt.Println(string(out))
	}
```

#### `song_test.go`

Keeping with the DRY prinicple, a common `gttTest()` function is used
for all the GTT tests. Since running a GTT test only varies by the
script to be run a common function makes a lot of sense. Give a script
file name a GTT UseCase is created follwed by creating a Runner with
the server information. The final step is to tell the Runner to Run.

``` golang
func gttTest(t *testing.T, filepath string) {
	uc, err := gtt.NewUseCase(filepath)
	if err != nil {
		t.Fatal(err.Error())
	}
	r := gtt.Runner{
		Server:   fmt.Sprintf("http://localhost:%d", testPort),
		Base:     "/graphql",
		Indent:   2,
		UseCases: []*gtt.UseCase{uc},
	}
	if testing.Verbose() {
		r.ShowComments = true
		r.ShowResponses = true
		r.ShowRequests = true
	}
	if err = r.Run(); err != nil {
		t.Fatal(err.Error())
	}
}
```

Each test is set up to run in a separate test function. This allows
sellecting individual tests from the command line with the go test
`-run` option. The tests themselves just call the `gttTest()` function
with the appropriate script file name.

``` golang
func TestTypes(t *testing.T) {
	gttTest(t, "gtt/types.json")
}
```

## Test Scripts

Script files are in a separate directory named `gtt`. More details on
the script file are in
[file_format.md](https://github.com/ohler55/graphql-test-tool/blob/master/file_format.md)
which can be referred to while we walk though the scripts.

Touching on a few of the features available in GTT lets look at a few
of the scripts.

 - **comments** If using the scripts to also document use cases
   comments are useful. Comments can be a single string value for a
   `comment` key as in
   [types.json](https://github.com/ohler55/graphql-test-tool/tree/master/example/gtt/types.json)
   or multiple lines in an array as in
   [artist_names_get.json](https://github.com/ohler55/graphql-test-tool/tree/master/example/gtt/artist_names_get.json).

 - TBD

 - sortBy - artist_names_get.json
 - remember - top.json
 - GET vs POST - artist_names_get.json and artist_names_post.json
 - leave off elements that don't matter - types.json


## Summary

While other packages were involved in managing the data behind the
GraphQL servers, [GGql](https://github.com/UHN/ggql) was key in
building the applications
quickly. [GTT](https://github.com/ohler55/graphql-test-tool) was key
to testing and providing examples to the teams writing the user
portions of the project.
