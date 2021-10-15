call main
exit


main:
	# 呼び出し元に戻れるように現状保存
	push bp
	cp sp bp
	# todo : 関数本体

	# 引数準備 : me
	sub_sp 1

	# ローカル変数に値を代入
	cp "john" [bp-1]

	push 18 # 引数(18)
	push [bp-1] # 引数(me)
	call say_hello # 関数呼び出し
	add_sp 2 # 呼び出し後の位置修正

	# 引数準備 : result
	sub_sp 1

	push 21 # 引数(21)
	push "tom" # 引数("tom")
	call say_hello # 関数呼び出し
	add_sp 2 # 呼び出し後の位置修正
	# 戻り値をローカル変数に格納
	cp reg_a [bp-2]

	echo [bp-2]

	# 現状復帰
	cp bp sp
	pop bp
	ret



say_hello:
	# 呼び出し元に戻れるように現状保存
	push bp
	cp sp bp

	# 引数をローカル変数として代入
	# your_name
	sub_sp 1
	cp [bp+2] [bp-3]
	# age
	sub_sp 1
	cp [bp+3] [bp-4]
	# todo : 関数本体

	echo [bp-3]

	echo [bp-4]


	# 戻り値を設定
	cp "success" reg_a

	# 現状復帰
	cp bp sp
	pop bp
	ret

