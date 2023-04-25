
# Oobabooga GUI Launcher

`smh...`

Not really tested anywhere ;-)

For original thing go to https://github.com/oobabooga/text-generation-webui

## Usage

Since Launcher supports multiple Obabooga installations (ie. different branches), `-home` argument is always required.

On Windows
```bash
dist\oobabooga.exe -install -home D:\oobabooga -- --model-dir D:\models --chat --auto-launch
```

works with git bash as well
```bash
dist/oobabooga.exe --home /d/oobabooga -- -h
```

That's basically it, GUI should open in your browser in chat mode.
See `oobabooga --help` and `oobabooga --home YOUR_HOME -- -h`

## Building

For Windows
```bash
go build -o dist/launcher.exe cmd/main.go
```

For Linux
```bash
go build -o dist/launcher cmd/main.go
```

## Building with Docker

On Windows
```bash
docker run --rm -ti -v %cd%:/go/src golang:1.20 bash -c "cd src; GOOS=windows go build -o dist/oobabooga.exe cmd/main.go"
```

On Linux
```bash
docker run --rm -ti -v $(pwd):/go/src golang:1.20 bash -c "cd src; go build -o dist/oobabooga cmd/main.go"
```

## TODO

Add support for zip bundles allowing to download whole environment as one big file (approx. 1GB) instead of pulling a bunch of small files.
