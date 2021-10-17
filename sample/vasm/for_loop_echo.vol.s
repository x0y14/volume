; print hello N times
; reg_a : 何回表示するか
; reg_b : 何回目のループか


echo "start"
push 5 ; 引数, 3回ループ
call for_i_in_range_X ; for関数を呼び出し
add_sp 1 ; spの位置修正
echo "finish"
exit ; 終了


; for loop, X: int
for_i_in_range_X:
    ; 準備
    push bp ; スタックにbpの位置保存
    cp sp bp ; spをbpに保存


    ; 引数をreg_aに保存
    cp [bp+2] reg_a

    ; init loop count
    cp 0 reg_b

    ; 本体は引数を取るので、引数を取らない内部ループを使用
    call i_loop


    ; 掃除
    cp bp sp ; bpからspにデータを返す
    pop bp ; stackからbpの位置を取り出す
    ret ; 関数から戻る


i_loop:
    echo "hello-loop"
    call i_increase_reg_b ; reg_b++
    cmp reg_a reg_b

    ; if not same
    jz i_loop

    ; if same
    ret



i_increase_reg_b:
    push reg_a ; reg_aのデータを保存します。
    cp 1 reg_a
    add reg_a reg_b
    pop reg_a ; reg_aのデータ取り出します。
    ret