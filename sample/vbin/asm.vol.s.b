call 3 
exit 
push bp 
cp sp bp 
sub_sp 1 
cp "yuhei" [bp-1] 
push 18 
push [bp-1] 
call 42 
add_sp 2 
sub_sp 1 
push 21 
push "tom" 
call 42 
add_sp 2 
cp reg_a [bp-2] 
echo [bp-2] 
cp bp sp 
pop bp 
ret 
push bp 
cp sp bp 
sub_sp 1 
cp [bp+2] [bp-3] 
sub_sp 1 
cp [bp+3] [bp-4] 
echo [bp-3] 
echo [bp-4] 
cp "success" reg_a 
cp bp sp 
pop bp 
ret 
