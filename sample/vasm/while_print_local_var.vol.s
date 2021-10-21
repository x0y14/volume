call main
exit


main:
    ; == 戻り用 ==
    push bp
    cp sp bp
    ; ===========

    ; == 関数本体 ==
    ; 条件式に必要なデータを、ローカル変数に代入

    sub 1 sp     ; ... : left
    cp 0 [bp-1]  ; ... : left
    sub 1 sp     ; 領域確保: right
    cp 5 [bp-2]  ; 代入: right

    ; whileに引数として渡す。
    push [bp-2] ; right
    push [bp-1] ; left
    call main_while_loop_01_entry
    add 2 sp     ; 引数の数、spを戻す
    ; =============

    ; == 呼び出し前の状態に復元 ==
    cp bp sp
    pop bp
    ret
    ; ========================


main_while_loop_01_entry:
    ; left opr right

    ; == 戻り用 ==
    push bp
    cp sp bp
    ; ===========
    jmp main_while_loop_01_conditional_expr

main_while_loop_01_conditional_expr:
    ; == 関数本体 ==
    cp [bp+2] reg_a ; left (0)
    cp 5 reg_b ; right(n)

    cmp reg_a reg_b ; !=
    ; もし、一致しなければ、本体へ飛ぶ。
    jz main_while_loop_01_body
    ; 終了へジャンプ
    jnz main_while_loop_01_end
    ; ============

main_while_loop_01_body:
    ; == 関数本体 ==
    push reg_a ; プリント引数
    call print ; プリント実行
    add 1 sp   ; プリント後始末

    ; nは、exprにおいて、rightに格納されていた。
    cp [bp+3] reg_a
    sub 1 reg_a ; n--をしたので、データをローカル変数に再代入
    cp reg_a [bp+3]
    jmp main_while_loop_01_conditional_expr
    ; =============

main_while_loop_01_end:
    ; == 呼び出し前の状態に復元 ==
    cp bp sp
    pop bp
    ret
    ; ========================


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