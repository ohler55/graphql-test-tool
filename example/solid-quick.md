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

The example GraphQL application is taken directly from the
[GGql](https://github.com/UHN/ggql) reflection example. GGql is the
fastest Go GraphQL server as well as being the easiest to use as shown
by this
[comparison](https://github.com/the-benchmarker/graphql-benchmarks/blob/develop/rates.md). The
schema in the example is:

``` graphql
type Query {
  artist(name: String!): Artist
  artists: [Artist]
  top: Song
}

type Mutation {
  like(artist: String!, song: String!): Song
  setLike(artist: String!, song: String!, count: Int!): Song
}

type Artist {
  name: String!
  songs: [Song]
  origin: [String]
}

type Song {
  name: String!
  artist: Artist
  duration: Int
  release: Date
  likes: Int
}

scalar Date
```

## Test Setup

 - choices
   - run as an app
     + framework for testing against any app, not just go
     - go test coverage tools don't work out of the box
   - run as part of test
     + output is displayed while running (debugging is easier)
     + go coverage tools work
     - go code needs to a cmd and package directory
 - gtt dir for scripts
 - with choice #1 a main_test.go start up app once and reuses
   - can run either or any server
 - what each test looks like

## Test Scripts

Script files are in a separate directory named `gtt`. More details on
the script file are in
[file_format.md](https://github.com/ohler55/graphql-test-tool/blob/master/file_format.md)
which can be referred to while we walk though the scripts.

 -

## Summary

While other packages were involved in managing the data behind the
GraphQL servers, [GGql](https://github.com/UHN/ggql) was key in
building the applications
quickly. [GTT](https://github.com/ohler55/graphql-test-tool) was key
to testing and providing examples to the teams writing the user
portions of the project.
