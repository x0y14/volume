echo "hello"
nop
echo 1
nop
echo 1.2
nop

cp "john" [sp-5]
cp [sp-5] reg_b
cp [sp-5] [bp-3]

echo [bp-3]
nop
echo reg_b
nop

exit