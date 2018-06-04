Break Out Of Hareka's Underground (Boohu) is a roguelike game mainly inspired from DCSS and its tavern, with some ideas from Brogue, but aiming for very short games, almost no character building, and a simplified inventory.

*Every year, the elders send someone to collect medicinal simella plants in the
Underground.  This year, the honor fell upon you, and so here you are.
According to the elders, deep in the Underground, magical stairs will lead you
back to your village.  Along the way, you will collect simellas, as well as
various items that will help you deal with monsters, which you may
fight or flee...*

![Boohu introduction screen](https://download.tuxfamily.org/boohu/intro-screen.png)

Screenshot
----------

[![Introduction Screeshot](https://download.tuxfamily.org/boohu/screenshot.png)](https://download.tuxfamily.org/boohu/index.html)

You can visit the [game's
page](https://download.tuxfamily.org/boohu/index.html)
for more informations, tips, screenshots and asciicasts. You will also be able
to play in the browser.

Install
-------

*main repo being migrated to tuxfamily*

You can download binaries on the [releases
page](https://github.com/anaseto/boohu/releases).

You can also build from source by following these steps:

+ Install the [go compiler](https://golang.org/).
+ Set `$GOPATH` variable (for example `export GOPATH=$HOME/go`).
+ Add `$GOPATH/bin` to your `$PATH` (for example `export PATH="$PATH:$GOPATH/bin"`).
+ Use the command `go get -u https://git.tuxfamily.org/boohu/boohu.git`.
  
The `boohu` command should now be available.

The only dependency outside of the go standard library is the lightweight
curses-like library [termbox-go](https://github.com/nsf/termbox-go), which is
installed automatically by the previous `go get` command.

*Portability note.* If you happen to experience input problems, try adding
option `--tags tcell` or `--tags ansi` to the `go get` command. The first will use
[tcell](https://github.com/gdamore/tcell) instead of termbox-go, and requires
cgo on some platforms, but is more portable. The second will work on POSIX
systems with a `stty` command.

Colors
------

If the default colors do not display nicely on your terminal emulator, you can
use the `-s` option: `boohu -s` to use the 16-color palette, which
will display nicely if the [solarized](http://ethanschoonover.com/solarized)
palette is used. Configurations are available for most terminal emulators, otherwise, colors may have to be configured manually to one's liking in
the terminal emulator options.
