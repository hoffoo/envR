envR
====

(experimental) Environment for R programming


ABOUT
-----

Manages an R process, letting you send snippets to it from your editor (rudimentary vim support included). A few other features are planned as well.


HOWTO
-----

```sh
go get github.com/hoffoo/envR
envR --no-save               # arguments passed to envR are all passed to R
```

This will start up R as you would usually run it. Subsequent envR invocations will be evaluated into the same workspace.

```sh
echo "1+4" | envR
```

You should see the expression evaluated in the original terminal.

In vim use :Rrun to send the current selection to R


BUGS
----

R crashes on error (not like in the standard REPL). Havent thought of how to fix it yet.


FUTURE
------

1. I'd like to be able to inspect the current workspace and get some kind of tagbar out for vim. This shouldn't be too difficult, there are other vim plugins that do it.
1. Find reliable way to forward signals to R
1. Simple syntax highlighting of output (such as argument count)


LICENCE
-------

MIT
