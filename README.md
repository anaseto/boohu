This is Break Out Of Hareka's Underground (Boohu), a roguelike game which takes
inspiration mainly from DCSS and its tavern, and some ideas from Brogue, but
aiming for very short games, almost no character building, and simplified
inventory management.

You somehow ended up in Hareka's Underground, and in order to leave, you have
to find some magical stairs deep in the Underground. Along the way, you will
collect treasures, as well as various helpful items. You will also encounter
monsters, fight against them, or run away from them when you can…

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

Survival Tips
-------------

Basic survival tips:

+ Position yourself to fight one ennemy at a time whenever possible.
+ Fight far enough from unknown terrain if you can: combat is noisy, more
  monsters will come if they hear it.
+ Use your potions, projectiles and rods. With experience you learn when you
  can spare them, and when you should use them, but dying with a potion of heal
  wounds in the inventory should never be a thing.
+ Avoid dead-ends if you do not have any means of escape, such as potions of
  teleportation, unless no better options are available.
+ Use *pillar dancing*: sometimes you can turn around a block several times to
  avoid being killed while replenishing your HP.
+ You do not have to kill every monster. You want, though, find as many items
  as you can, but survival comes first.
