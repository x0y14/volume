call main
exit

main:
    ; == 戻り用 ==
    push bp
    cp sp bp
    ; ===========


    ; == ローカル変数定義: n ==
    sub 1 sp
    ; ローカル変数に値を代入
    cp 5 [bp-1]
    ; ======================


    ; == 呼び出し前の状態に復元 ==
    cp bp sp
    pop bp
    ret
    ; ========================