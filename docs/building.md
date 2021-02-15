# How To Build from source

## macOS

1. Install [Go](https://golang.org/) via [brew](https://brew.sh/)

```bash
$ brew install go  
$ go version
go version go1.10.2 darwin/amd64
```

1. Install Go cross compilation tool - [gox](https://github.com/mitchellh/gox) via [brew](https://brew.sh/)

```bash
$ brew install gox  
$ brew info gox  
gox: stable 0.4.0 (bottled)
```

1. Check [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable)

```bash
$ printenv GOPATH  
/Users/<user>/go
```

1. If [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable) not defined, add it to bash profile

```bash
$ echo 'export GOPATH="$HOME/go"' >> ~/.bash_profile  
$ source ~/.bash_profile  
$ printenv GOPATH  
/Users/<user>/go
```

1. Get goloc package

```bash
go get github.com/s0nerik/goloc
```

1. Move to package folder

```bash
cd <GOPATH>/src/github.com/s0nerik/goloc
```

1. Run build script and wait

```bash
$ sh build.sh  
Number of parallel builds: 3
-->   windows/amd64: github.com/s0nerik/goloc
-->     linux/amd64: github.com/s0nerik/goloc
-->    darwin/amd64: github.com/s0nerik/goloc
adding: darwin_amd64 (deflated 68%)
adding: linux_amd64 (deflated 68%)
adding: windows_amd64.exe (deflated 68%)
```

1. Move to output folder and unzip binary

```bash
$ cd out  
$ unzip goloc.zip  
inflating: darwin_amd64
inflating: linux_amd64
inflating: windows_amd64.exe
```
