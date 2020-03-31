To run:
```
go run ./cmd help
go run ./cmd sourcefile targetfile
```

To build binary:
```
go build -o ./bin/fuzzyelem ./cmd
```

Binary usage:
```
./bin/fuzzyelem help
./bin/fuzzyelem --id make-everything-ok-button ./samples/sample-0-origin.html ./samples/sample-4-the-mash.html
```