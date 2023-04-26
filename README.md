
# Text Generation WEBUI Launcher

Portable [text-generation-webui](https://github.com/oobabooga/text-generation-webui) installation. Independent of already installed Python or Conda environments. Self-contained setup. Single binary entrypoint.

TBH no idea if anything I described above is actually true. This project was put together in few minutes, so... Not really tested anywhere ;-)


## Download

https://github.com/r3dsh/text-generation-webui-launcher/releases

## Usage

Since Launcher supports multiple [text-generation-webui](https://github.com/oobabooga/text-generation-webui) installations (ie. different branches), `-home` argument is always required.

#### On Windows
```bash
dist\text-generation-webui-launcher.exe -install -home D:\oobabooga -- --model-dir D:\models --chat --auto-launch
```

That's basically it, after installation is done, GUI should open in your browser in chat mode. Other examples:

#### Git bash as well
```bash
dist/text-generation-webui-launcher.exe --home /d/oobabooga -- -h
```

#### And WSL (linux binary)
```bash
└─$ ./text-generation-webui-launcher -home /mnt/d/oobabooga
2023/04/26 06:33:11 Starting text-generation-webui from /mnt/d/oobabooga/text-generation-webui-main
2023/04/26 06:33:11    spawning server.py
```

All from one directory. Python packages and environment variables are managed by launcher.

See `text-generation-webui-launcher --help` and `text-generation-webui-launcher --home YOUR_HOME -- -h`

## Building

For Windows
```bash
GOOS=windows go build -o dist/text-generation-webui-launcher.exe cmd/main.go
```

For Linux
```bash
go build -o dist/text-generation-webui-launcher cmd/main.go
```

For MacOS
```bash
GOOS=darwin go build -o dist/text-generation-webui-launcher cmd/main.go
```

## Building with Docker

On Windows
```bash
docker run --rm -ti -v %cd%:/go/src golang:1.20 bash -c "cd src; GOOS=windows go build -o dist/text-generation-webui-launcher.exe cmd/main.go"
```

On Linux
```bash
docker run --rm -ti -v $(pwd):/go/src golang:1.20 bash -c "cd src; go build -o dist/text-generation-webui-launcher cmd/main.go"
```

## TODO

- Fetch GIT binaries (launcher.go line 122)
- Add support for zip bundles allowing to download whole environment as one big file (approx. 1GB) instead of pulling a bunch of small files.
- "Select Directory" OS native window for non CLI usage instead of -home directory + config file? Something really simple, no additional UIs.
- Actually test stuff...