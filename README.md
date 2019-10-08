# mssfmt

Tool for T-SQL (Microsoft SQL Server) automatic code format.

## Getting Started

This tool takes T-SQL script(s) (single or whole catalog) and rewrite it similar
to Go's `go fmt` which is a tool for automatically formats Go code.
This `mssfmt` tool is meant to be as simple as `go fmt`. In opposite to `gm fmt`
`mssfmt` doesn't create abstract syntax tree of T-SQL it rather performs series
of regexps and transformation to prettify the source file. 

We are in early stage of developing this tool, please do not use this while serious
developing yet.

`mssfmt` is developed in `Go` therefore it is cross-platform and can be used
outside of Windows (in case you have SQL Server on your Linux server).

### Examples

On Windows:

```
mssfmt.exe InputTSqlScript.sql
```

### Installing

At this point binary isn't prepared and distributed - it would be after the first
realise. In case you want check it out earlier you have to got
[Go](https://golang.org/dl) installed.

```
git clone https://github.com/DSkrzypiec/mssfmt
cd mssfmt
go build
```

## Running the tests

```
go test ./... -cover
```


## Built With

* [Go](https://golang.org/) - go version go1.13 


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
