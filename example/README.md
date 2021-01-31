# Simple Singing Server

Two servers are available. One writtent in Go and one written in Ruby.
The Ruby implementation requires the
[Agoo]((https://github.com/ohler55/agoo) gem which can be installed
with this command:

```
gem install agoo
```

To start the Ruby server type the following:

```
ruby song.rb
```

The server has verbosity set pretty high so that it is possible to see
what is occurring. From a browser a simple query can be used to verify
the server is running.

`localhost:6464/graphql?query={artist(name:"Fazerdaze"){name}}&indent=2`

Once started the GraphQL-Test-Tool can be used with the examples files.

The Go server can be started with the command:

```
go run main.go
```

For more details read the [solid-quick.md](solid-quick.md) article.
