call 3
exit 
push bp
cp sp bp
sub 1 sp
cp 0 [bp-1]
sub 1 sp
cp 5 [bp-2]
push [bp-2]
push [bp-1]
call 35
add 2 sp
cp bp sp
pop bp
ret 
push bp
cp sp bp
jmp 42
cp [bp+2] reg_a
cp 5 reg_b
cmp reg_a reg_b
jz 55
jnz 73
push reg_a
call 79
add 1 sp
cp [bp+3] reg_a
sub 1 reg_a
cp reg_a [bp+3]
jmp 42
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