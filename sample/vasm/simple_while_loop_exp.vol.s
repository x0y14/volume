nop
call main
nop
exit


main:
    ;; store
    push bp
    cp sp bp

    ; ローカル変数, 3
    sub 1 sp
    cp 3 [bp-1]

    ; ローカル変数, 0
    sub 1 sp
    cp 0 [bp-2]

    ; 表示
    push [bp-1] ; 3
    push [bp-2] ; 0
    call main_while_loop_entry
    add 2 sp

    ;; restore
    cp bp sp
    pop bp
    ret



main_while_loop_entry:
    ;; store
    push bp
    cp sp bp

    ; goto conditional expr
    jmp main_while_loop_cond_expr


main_while_loop_cond_expr:
    ; reserve arguments
    cp [bp+2] reg_a ; 0
    cp [bp+3] reg_b ; 3

    ; (reg_a < reg_b)
    lt reg_a reg_b
    jnz  main_while_loop_body
    jz main_while_loop_end


main_while_loop_body:
    push [bp+2]
    call print
    add 1 sp

    ; decrement
    cp [bp+2] reg_a
    add 1 reg_a
    cp reg_a [bp+2]

    ; back to cond expr
    jmp main_while_loop_cond_expr


main_while_loop_end:
    ;; restore
    cp bp sp
    pop bp
    ret



; link: stdio.vol.s
; (linkerまだないので、手動リンクです。)
print:
    ; == 戻り用 ==
    push bp
    cp sp bp
    ; ===========

    ; == 関数本体 ==
    ; 表示する内容を引数にとり、ローカル変数に。
    sub 1 sp
    ; why 3? => 知らん過去のスクリプト参照した。
    cp [bp+2] [bp-3]
    ; 表示
    echo [bp-3]
    ; =============

    ; == 呼び出し前の状態に復元 ==
    cp bp sp
    pop bp
    ret
    ; ========================