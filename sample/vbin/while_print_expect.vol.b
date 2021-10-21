call 3
exit 
push bp
cp sp bp
cp 5 reg_a
call 19
cp bp sp
pop bp
ret 
push bp
cp sp bp
jump 26
cmp reg_a 0
jz 33
jump 45
push reg_a
call 51
; print
add 1 sp
sub 1 reg_a
jump 26
cp bp sp
pop bp
ret 
push bp
cp sp bp
sub 1 sp
cp [bp+2] [bp-3]
echo [bp-3]
cp bp sp
pop bp
ret 