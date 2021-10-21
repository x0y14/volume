call 3; main
exit 
push bp
cp sp bp
cp 5 reg_a
call 19; main_while_loop_01_entry
cp bp sp
pop bp
ret 
push bp
cp sp bp
jump 26; main_while_loop_01_conditional_expr
cmp reg_a 0
jz 33; main_while_loop_01_body
jump 45; main_while_loop_01_end
push reg_a
call 51; print
add 1 sp
sub 1 reg_a
jump 26; main_while_loop_01_conditional_expr
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