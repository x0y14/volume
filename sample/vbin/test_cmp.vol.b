; eq
cp "yes" reg_a
cp "yes" [sp-5]
cmp reg_a [sp-5]
echo zf
; 1

; 0 < 1
lt 0 1
echo zf
; 1
cp 0 reg_a
cp 1 reg_b
lt reg_b reg_a
echo zf
; 0

; 0 > 1
gt 0 1
echo zf
; 0
cp 1 reg_a
cp 0 reg_b
gt reg_a reg_b
echo zf
; 1


exit