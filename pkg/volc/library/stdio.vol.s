print:
    ; arguments
    ; 0: text

	; == prepare to back ==
	push bp
	cp sp bp
	; =====================

    ; == function arguments ==
	; text
	sub_sp 1
	cp [bp+2] [bp-3]
	; ========================

    ; == function body ==
	echo [bp-3]
	; ===================

	; == do back ==
	cp bp sp
	pop bp
	ret
	; =============