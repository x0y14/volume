call 3
exit 
push bp
cp sp bp
sub 1 sp
cp 5 [bp-1]
push [bp-1]
call 27
add 1 sp
cp bp sp
pop bp
ret 
push bp
cp sp bp
jump 34
cp [bp+2] reg_a
cmp reg_a 0
jz 44
jump 59
push reg_a
call 65
add 1 sp
sub 1 reg_a
cp reg_a [bp+2]
jump 34
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