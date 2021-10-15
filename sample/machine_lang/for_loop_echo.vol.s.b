echo "start" 
push 5 
call 11 
add_sp 1 
echo "finish" 
exit 
push bp 
cp sp bp 
cp [bp+2] reg_a 
cp 0 reg_b 
call 30 
cp bp sp 
pop bp 
ret 
echo "hello-loop" 
call 40 
cmp reg_a reg_b 
jz 30 
ret 
push reg_a 
cp 1 reg_a 
add reg_a reg_b
pop reg_a 
ret 
