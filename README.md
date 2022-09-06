# FileExtractor

Tool for extracting files with given extensions into a separate output directory.

# Usage
```bash
./fileextractor -src ./sourcedirectory -dest ./destination -exts png jpg

# For help use
./filextractor -h
```

# Build binary or run directly

## Build binary
```bash
cd ./fileextractor
go build 
```

## Run directly
```bash
# Run directly with go
cd ./fileextractor
go run fileexctractor.go -exts png jpg
```

# Arguments
`-exts`: list of extensions to filter for

`-src` (default: `.`): specifies the source directory 

`-out` (default: `./output`): specifies the output directory 


## Examples
### Filter by single extension in given directory
```bash
./fileextractor -src ./mydir -exts pdf
```

### Filter by multiple extensions
```bash
./fileextractor -exts jpg png gif
```