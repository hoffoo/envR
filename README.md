envR
====

Environment for R programming


ABOUT
=====

Manages an R process, letting you send snippets to it from your editor (rudementary vim support included). A few other features are planned as well.


HOWTO
=====

```sh
go get github.com/hoffoo/envR
envR
```

This will start up R as you would usually run it. Subsequent envR invocations will be evaluated in that terminal, and input will be piped to R. Example:

```sh
echo "1+4" | envR
```

You should see the expression evaluated in the original terminal.


VIM
===

Use :Rrun to send the current selection to R


BUGS
====

I havent gotten the normal R stdin REPL to work properly, since the process takes over it. You can quit with ctrl+c.


FUTURE
======

Id like to be able to inspect the current workspace and get some kind of tagbar out for vim. This shouldn't be too difficult.


LICENCE
=======

MIT
