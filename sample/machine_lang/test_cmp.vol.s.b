cp "yes" reg_a
cp "yes" [sp-5]

cmp reg_a [sp-5]
nop
cmp reg_a "[sp-5]"
nop
exit