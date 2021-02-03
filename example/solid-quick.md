# Building Solid Go GraphQL Applications Quickly

After numerous successful deployments of Go GraphQL applications a
repeatable building methodology has emerged. Two of the most
significant contributing factors to the successful development and
deployment of those projects are a good GraphQL server package to
build the application and suitable black box testing tool.

## Black Box Testing

It would be hard to argue that application level black box testing is
not the gold standard for testing. Sure, unit testing is important but
the final gate before release or deployment really should be tests
that validate behavior as the end users will see it.

GraphQL servers are no exception to the black box testing rule. A test
suite that exercises an GraphQL application using the same HTTP API
that end users will use make it less likely there will be surprises
after deployment. A black box test suite is also ideal for continuous
integration tests on merges during development.

With the proper tooling, testing through the HTTP API of a GraphQL
application is often easier than trying to piece together regression
tests that only exercise internal APIs. An effective tool should be
able to play the role of a
user. [GraphQL-Test-Tool](https://github.com/ohler55/graphql-test-tool)
is such a tool. GTT for short,
[GraphQL-Test-Tool](https://github.com/ohler55/graphql-test-tool) is a
script driven test tool for testing GraphQL applications. I addition
to providing repeatable regression tests and CI, GTT scripts end up
being great examples for end users.

GTT is written in Go but can be used for testing any GraphQL server
that has an HTTP API. The [example](.) this article describes
demonstrates the use of GTT as a test tool for a Go application test
as well as a Ruby server that implements the same schema. This article
focuses on the Go test setup.

## Application

The GraphQL application is this example is taken from the
[GGql](https://github.com/UHN/ggql) reflection example. A few changes
to the example so that more of the GTT features could be explained but
basically it is the same. GGql is the fastest Go GraphQL server as
well as the easiest to use as shown by this
[comparison](https://github.com/the-benchmarker/graphql-benchmarks/blob/develop/rates.md).
The
[README.md](https://github.com/UHN/ggql/tree/master/examples/reflection/README.md)
for the GGql example explains the basics of the application so that
will not be duplicated here.

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

Another addition was of a `setLike()` mutation field. Since we are
using the GGql reflection approach to resolving fields all that has to
be done is add the Mutation function and GGql takes care of the rest.

``` golang
func (m *Mutation) SetLike(artist, song string, count int64) *Song {
	if a := m.query.Artist(artist); a != nil {
		if s := a.Song(song); s != nil {
			s.Likes = int(count)
			return s
		}
	}
	return nil
}
```

The GGql example did not include support for passing in variable nor
the operation name as query parameters in the URL. That functionality
was also added with these lines in the `handleGraphQL()` function.

``` golang
	var vars map[string]interface{}
	if variables, err := oj.ParseString(req.URL.Query().Get("variables")); err == nil {
		vars, _ = variables.(map[string]interface{})
	}
	op := req.URL.Query().Get("operationName")
```

With the application ready for testing lets move on to the test setup.

## Test Setup

There are two choices when setting up the tests for an
application. One is to use the full application for a true black box
test and the other is to split the application into a command (cmd)
portion and a package portion so that the package code can be tested
in the same process space as the test code. There are advantages and
disadvantages to each.

### True Black Box

A true black box approach runs the application completely separate
from the test code. This had the advantage of being able to test a
server implemented in any language. The disadvantage is that the Go
test coverage tools will not work. Running as a separate process also
means debugging print statements are a little more difficult to
display.

### Embedded Tests

In order to run the application in the same code space as the test
code the application needs to be callable from the test code. To do
that a package with the application in it needs to be imported. Thats
not difficult to do. Just create a `cmd` directory and put a light
weight `main()` function in an application directory of the `cmd`
directory that calls the package where all the rest of the code
resides. Go coverage tools then work and debug print statements show
up as the tests are running. The downside to an embedded configuration
is that the server has to be written in Go and separated in a cmd and
package directory. Of course writing the application in Go isn't
really much of a downside.

### Implementation

The true black box testing approach is more applicable to a wider
audience and has a certain testing purity to it so that is the
approach described here. Jumping right in, the test scripts are placed
in a `gtt` subdirectory to keep files organized. They can be placed
anywhere though.

While not necessary, the application will be started just once and
then each test will be executed against the running app. The code to
run the application is in `main_test.go` while the test functions for
running the individual tests are in `song_test.go`.

#### `main_test.go`

Since the tests can not be run without first starting the application
it will be necessary to create a `TestMain()` function. `TestMain()`
begins with setting up a flag to determine if the Go or Ruby version
of the application will be called. Next the `run()` function is
called. The `run()` function will do most of the work and return an
error if anything goes wrong with the setup or if the tests fail.

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

The `run()` function first finds a free port to run the application
on. This avoids accidental collisions with other servers or with
previously tests that may be taking longer to shutdown than expected.

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

With a free port identified, the application is started. Since we want
to collect the output from the `exec.Cmd` we grab the application
stdout before calling `start()`.

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
attempting to write to stdout. By reading concurrently the
`ioutil.ReadAll()` buffer will expand and collect all the
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

The tests shouldn't be run until the application has started or else
the tests will fail to connect to the application. A simple sleep
could be used then that forces a slowdown of the test. It's better to
continue immediately once the application is up and accepting
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

If there were no errors connecting to the application the tests can be
run.

``` golang
	if err == nil && 0 != m.Run() {
		err = fmt.Errorf("tests failed")
	}
}
```

After the tests finish it's time to clean up by killing the
application and printing the application output. Note the wait on the
`done` chan before proceeding to printing to make sure all the output
has been collected.

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

Individual tests all follow the same pattern so keeping with the DRY
principle, a common `gttTest()` function is used for all the GTT
tests. The only variable in the GTT tests is the script file to be
executed so that is passed as an argument to the `gttTest()` function.
A GTT UseCase is created with the filepath followed by the creation of
a `gtt.Runner` with the server information. The final step is to tell
the Runner to Run.

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
selecting individual tests from the command line with the go test
`-run` option. The tests themselves just call the `gttTest()` function
with the appropriate script file name.

``` golang
func TestTypes(t *testing.T) {
	gttTest(t, "gtt/types.json")
}
```

All the interactions with the server are in the script file.

## Test Scripts

Script files are in a separate directory named `gtt`. More details
about the structure of the script files can be found in
[file_format.md](https://github.com/ohler55/graphql-test-tool/blob/master/file_format.md).

Lets look at a few of the features available in GTT by touching on a
few of of the scripts.

 - **comments** - Comments in the script file are useful if using the
   scripts to also document use cases. Comments can be a single string
   value for a `comment` key as in [types.json](gtt/types.json) or
   multiple lines in an array as in
   [artist_names_get.json](gtt/artist_names_get.json).

 - **sortBy** - Order is often important in arrays yet at times order
   may be be random depending on the data source that returns a list of
   items. The `sortBy` option will sort output before comparing to the
   expected. [artist_names_get.json](gtt/artist_names_get.json) is an
   example of using the `sortBy` option.

 - **method** - Switching between HTTP methods, either GET or POST is
   as simple as either providing a `content` element or not as seen in
   [artist_names_get.json](gtt/artist_names_get.json) and
   [artist_names_post.json](gtt/artist_names_post.json).

 - **remember** - Sometimes the result of one step is needed in a
   subsequent step. The `remember` and `vars` elements provide that
   functionality in [top.json](gtt/top.sen).

 - **lazy** - GTT is extremely tolerant of script format and supports
   the [SEN](https://github.com/ohler55/ojg/blob/develop/sen.md)
   (Simple Encoding Notation). The relaxed format is used in
   [top.json](gtt/top.sen) but it also allow for JSON with errors such
   as missing commas or extra commas.

## Summary

With a full set of use cases a GraphQL application can be deployed
with confidence. Having the ability to run regression tests after
making changes saves time especially when a hot fix is needed. A much
appreciated feature of scripted tests is that users have examples they
can use when putting together we front ends.
