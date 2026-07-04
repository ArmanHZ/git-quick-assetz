# Git Quick AssetZ (quaz)
A TUI [(tview)](https://github.com/rivo/tview) tool written in Go for downloading assets from git releases with ease.

# Why?
I needed a tool with a user-friendly UI to quickly hand pick and download specific assets from
releases of tool repository when working on CTF/HTB challenges.

___Note: `quaz` only displays the non-default assets for each release. In other words, source `zip` and `tar.gz` files are not shown. This is intended behavior, because downloading them are no different to cloning the repository.___

# Build and run
Install using `go` ($GOPATH must be in your environment variables in order to use `quaz` from everywhere):

```console
go install github.com/ArmanHZ/git-quick-assetz/cmd/quaz@latest
```

To build and run, you can use the `build_and_run.sh` or manually as follows:

```console
git clone https://github.com/ArmanHZ/git-quick-assetz.git
go build -o ./quaz ./cmd/quaz/quaz.go
./quaz
```

If you want, you can move the `quaz` to somewhere that is in your path and use it
from anywhere you want.

# Usage
The movement keys are `Tab` and `Shift+Tab`.

The action key is `Enter`.

# Screenshots and functionalities
With the default color scheme of `Windows Terminal`, the TUI app looks like this:

![grd-1](./img/grd-1.png)

You can enter/paste a valid `GitHub` URL in the text box and press `enter` to pull
all available releases:

![grd-2](./img/grd-2.png)

You can navigate to the release you want, expand it and select the files you want:

![grd-3](./img/grd-3.png)

After that, you can press the `Download Assets` button to see the assets that will
be downloaded and as well as the ability to define a directory where they will be
downloaded to.

![grd-4](./img/grd-4.png)

By default, your current working directory will be shown in the `Save location`
input field.
