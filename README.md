
# Text Generation WEBUI Launcher

Portable https://github.com/oobabooga/text-generation-webui installation. Independent of already installed Python or Conda environments. Self-contained setup. Single binary entrypoint.

TBH no idea if anything I described above is actually true. This project was put together in few minutes, so... Not really tested anywhere ;-)


## Download

https://github.com/r3dsh/text-generation-webui-launcher/releases

## Usage

Since Launcher supports multiple text-generation-webui installations (ie. different branches), `-home` argument is always required.

On Windows
```bash
dist\launcher.exe -install -home D:\oobabooga -- --model-dir D:\models --chat --auto-launch
```

works with git bash as well
```bash
dist/launcher.exe --home /d/oobabooga -- -h
```

That's basically it, GUI should open in your browser in chat mode.
See `launcher --help` and `launcher --home YOUR_HOME -- -h`

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
docker run --rm -ti -v %cd%:/go/src golang:1.20 bash -c "cd src; GOOS=windows go build -o dist/launcher.exe cmd/main.go"
```

On Linux
```bash
docker run --rm -ti -v $(pwd):/go/src golang:1.20 bash -c "cd src; go build -o dist/launcher cmd/main.go"
```

## TODO

Add support for zip bundles allowing to download whole environment as one big file (approx. 1GB) instead of pulling a bunch of small files.
