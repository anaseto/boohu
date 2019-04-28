**This sealth branch is an experiment for a stealth focused variant of Boohu.
In case it actually becomes something, it would be a different game
(Harmonist).**

Harmonist is a stealth coffee-break roguelike game.  The game has a heavy focus
on tactical positioning, light and noise mechanisms, making use of various
terrain types and cones of view for monsters.  Aiming for a replayable
streamlined experience, the game avoids complex inventory management and
character building, relying on items and player adaptability for character
progression.

*Your friend Shaedra got captured by nasty people from the Dayoriah Clan while
she was trying to retrieve a powerful magara artifact that was stolen from the
great magara-specialist Marevor Helith. As a gawalt monkey, you don't
understand much why people complicate so much their lives caring about
artifacts and the like, but one thing is clear: you have to rescue your friend,
somewhere to be found in this Underground area controlled by the Dayoriah Clan.
If what you heard the guards say is true, Shaedra's imprisoned on the eighth
floor. You are small and have good night vision, so you hope the infiltration
will go smoothly...*

TODO: update the rest of the README.md when there is a website for the new game.

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

    go get -u git.tuxfamily.org/boohu/boohu.git --tags tk

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
