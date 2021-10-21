nop 
call 5
nop 
exit 
push bp
cp sp bp
sub 1 sp
cp 3 [bp-1]
sub 1 sp
cp 0 [bp-2]
push [bp-1]
push [bp-2]
call 37
add 2 sp
cp bp sp
pop bp
ret 
push bp
cp sp bp
jmp 44
cp [bp+2] reg_a
cp [bp+3] reg_b
lt reg_a reg_b
jnz 57
jz 75
push [bp+2]
call 81
add 1 sp
cp [bp+2] reg_a
add 1 reg_a
cp reg_a [bp+2]
jmp 44
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