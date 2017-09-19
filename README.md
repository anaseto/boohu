This is Break Out Of Hareka's Underground (short BOOHU), a roguelike game which
takes inspiration mainly from DCSS and its tavern, and some ideas from Brogue,
but aiming for very short games, almost no character building, and simplified
inventory management.

You somehow ended up in Hareka's Underground, and in order to leave, you have
to find some magical stairs deep in the Underground. Along the way, you will
collect treasures, as well as various items to help yourself. You will also
encounter monsters, fight against them, or run away from them when you can…

It is a work in progress, but is already a quite complete game.

![boohu intro screen](https://raw.githubusercontent.com/anaseto/boohu/master/img/intro-screen.png)

Install
-------

+ Install the [go compiler](https://golang.org/).
+ Set `$GOPATH` variable (for example `export GOPATH=$HOME/go`).
+ Add `$GOPATH/bin` to your `$PATH` (for example `export PATH="$PATH:$GOPATH/bin"`).
+ Use the command `go get -u github.com/anaseto/boohu`.
  
The `boohu` command should now be available.

The only dependency outside of the go standard library is the lightweight
curses-like library [termbox-go](https://github.com/nsf/termbox-go), which is
installed automatically by the previous `go get` command.

Colors
------

If the default colors do not display nicely on your terminal emulator, you can
use the `-s` option: `boohu -s`. The game then uses the 16-color palette, and
will display nicely if the [solarized](http://ethanschoonover.com/solarized)
palette, with configurations available easily for most terminal emulators, is
used.  Otherwise, colors may have to be configured manually to one's liking in
the terminal emulator options.
