### NOP
`nop`  
何もしない

### SET (Deprecated)
`set src dst` => `dst = src`

### ADD
`add a1 a2` -> `a2 += a1`  
a1: (int, float)  
a2: (int, float)  

### SUB
`sub a1 a2` -> `a2 -= a1`  
a1: (int, float)  
a2: (int, float)

### CMP
`cmp a1 a2`  
a1 == a2 -> `zf = 1`  
a1 != a2 -> `zf = 0`

### SJOIN
`sjoin a1 a2` => `a2 += a1`  
a1: (string)  
a2: (string)

### JUMP

### JZ

### JNZ

### CALL

### RET

### CP

### PUSH

### POP

### ADDsp

### SUBsp

### ECHO

### EXIT