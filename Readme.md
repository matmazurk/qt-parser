# qt-parser

### Simple QuickTime parser

Includes `parser` package for API and executable which reads dimensions from video tracks and sampling frequency from audio tracks

## Installation
```bash
go install github.com/matmazurk/qt-parser@latest
```

## Running
```bash
qt-parser path/to/file
```

## Running without installation
```bash
git clone github.com/matmazurk/qt-parser
cd qt-parser
go run . path/to/file
```

## Running tests
```bash
go test ./...
```
