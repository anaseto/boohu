Break Out Of Hareka's Underground (Boohu) is a roguelike game mainly inspired
from DCSS and its tavern, with some ideas from Brogue, but aiming for very
short games, almost no character building, and a simplified inventory.

*Every year, the elders send someone to collect medicinal simella plants in the
Underground.  This year, the honor fell upon you, and so here you are.
According to the elders, deep in the Underground, a magical monolith will lead you
back to your village.  Along the way, you will collect simellas, as well as
various items that will help you deal with monsters, which you may
fight or flee...*

![Boohu introduction screen](https://download.tuxfamily.org/boohu/intro-screen-tiles.png)

Screenshot and Website
----------------------

[![Introduction Screeshot](https://download.tuxfamily.org/boohu/screenshot.png)](https://download.tuxfamily.org/boohu/index.html)

You can visit the [game's
website](https://download.tuxfamily.org/boohu/index.html)
for more informations, tips, screenshots and asciicasts. You will also be able
to play in the browser and download pre-built binaries for the latest release.

Install from Sources
--------------------

In all cases, you need first to perform the following preliminaries:

+ Install the [go compiler](https://golang.org/).
+ Set `$GOPATH` variable (for example `export GOPATH=$HOME/go`, the default
  value in recent Go versions).
+ Add `$GOPATH/bin` to your `$PATH` (for example `export PATH="$PATH:$GOPATH/bin"`).

### ASCII

You can build a native ASCII version from source by using this command:

+ `go get -u git.tuxfamily.org/boohu/boohu.git`.
  
The `boohu` command should now be available (you may have to rename it to
remove the `.git` suffix).

The only dependency outside of the go standard library is the lightweight
curses-like library [termbox-go](https://github.com/nsf/termbox-go), which is
installed automatically by the previous `go get` command.

*Portability note.* If you happen to experience input problems, try adding
option `--tags tcell` or `--tags ansi` to the `go get` command. The first will use
[tcell](https://github.com/gdamore/tcell) instead of termbox-go, which is more
portable (works on OpenBSD). The second will work on POSIX systems with a
`stty` command.

### Tiles

You can build a graphical version depending on Tcl/Tk (8.6) using this command:

    go get -u --tags tk git.tuxfamily.org/boohu/boohu.git

This will install the [gothic](https://github.com/nsf/gothic) Go bindings for
Tcl/Tk. You need to install Tcl/Tk first.

With Go 1.11 or later, you can also build the WebAssembly version with:

    GOOS=js GOARCH=wasm go build --tags js -o boohu.wasm

You can then play by serving a directory containing the wasm file via http. The
directory should contain some other files that you can find in the main
website instance.

Colors
------

If the default colors do not display nicely on your terminal emulator, you can
use the `-s` option: `boohu -s` to use the 16-color palette, which
will display nicely if the [solarized](http://ethanschoonover.com/solarized)
palette is used. Configurations are available for most terminal emulators,
otherwise, colors may have to be configured manually to one's liking in
the terminal emulator options.

Documentation
-------------

See the man page boohu(6) for more information on command line options and use
of the replay file. For example:

    boohu -r _

launches an auto-replay of your last game.
