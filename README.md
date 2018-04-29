This is Break Out Of Hareka's Underground (Boohu), a roguelike game which takes
inspiration mainly from DCSS and its tavern, and some ideas from Brogue, but
aiming for very short games, almost no character building, and simplified
inventory management.

Every year, your village sends someone to collect medicinal simella plants in
the Underground.  This year, the duty fell upon you, and so here you are. Your
heart is teared between your will to be as helpful as possible to your village
and your will to make it out alive.  Deep in the Underground, some magical
stairs will lead you back to your village. Along the way, you will collect
simellas, as well as various helpful items. You will also encounter monsters,
fight against them, or run away from them when you can…

![boohu intro screen](https://raw.githubusercontent.com/anaseto/boohu/master/img/intro-screen.png)

Screenshot
----------

![Dragons](https://raw.githubusercontent.com/anaseto/boohu/master/img/dragons.png)

There are also some [asciinema screencasts](https://asciinema.org/~anaseto).

Install
-------

You can found binaries in the [releases
page](https://github.com/anaseto/boohu/releases).

You can also build from sources by following these steps:

+ Install the [go compiler](https://golang.org/).
+ Set `$GOPATH` variable (for example `export GOPATH=$HOME/go`).
+ Add `$GOPATH/bin` to your `$PATH` (for example `export PATH="$PATH:$GOPATH/bin"`).
+ Use the command `go get -u github.com/anaseto/boohu`.
  
The `boohu` command should now be available.

The only dependency outside of the go standard library is the lightweight
curses-like library [termbox-go](https://github.com/nsf/termbox-go), which is
installed automatically by the previous `go get` command.

*Portability note.* If you happen to experience input problems, try adding
option `--tags tcell` to the `go get` command, which will use
[tcell](https://github.com/gdamore/tcell) instead of termbox-go. It requires
cgo, but is more portable.

Colors
------

If the default colors do not display nicely on your terminal emulator, you can
use the `-s` option: `boohu -s`. The game then uses the 16-color palette, and
will display nicely if the [solarized](http://ethanschoonover.com/solarized)
palette, with configurations available easily for most terminal emulators, is
used.  Otherwise, colors may have to be configured manually to one's liking in
the terminal emulator options.

Basic Survival Tips
-------------

+ Position yourself to fight one enemy at a time whenever possible.
+ Fight far enough from unknown terrain if you can: combat is noisy, more
  monsters will come if they hear it.
+ Use your potions, projectiles and rods. With experience you learn when you
  can spare them, and when you should use them, but dying with a potion of heal
  wounds in the inventory should never be a thing.
+ Avoid dead-ends if you do not have any means of escape, such as potions of
  teleportation, unless no better options are available.
+ Use *pillar dancing*: sometimes you can turn around a block several times to
  avoid being killed while replenishing your HP.
+ You do not have to kill every monster. You want, though, to find as many items
  as you can, but survival comes first.
+ Use doors and dense foliage to break line of sight with monsters.
