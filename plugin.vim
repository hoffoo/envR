function! Rrun()
    call system("envR", join(getline(a:firstline, a:lastline), "\n"))
endfunction

com! -range=% -nargs=0 Rrun :<line1>,<line2>call Rrun()
