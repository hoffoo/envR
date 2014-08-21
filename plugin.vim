function! Rrun() range
    call system("envR", join(getline(a:firstline, a:lastline), "\n"))
endfunction

com! -range=% -nargs=0 Rrun :'<,'>call Rrun()
