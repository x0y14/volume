## 対応状況
*関数関係*
- [ ] 定義
- [x] 仮引数
- [x] 戻り値
- [ ] 本体
- [ ] 呼び出し
- [ ] 呼び出しに使用される引数

*本文*
- [ ] コンテンツ

*インポート*
- [x] import

*リターン*
- [ ] return

*条件分岐*
- [ ] if
- [ ] elif
- [ ] else
- [ ] 分岐の本体

*変数*
- [ ] 定義&代入
- [ ] 代入
- [ ] データ

*式*
- [ ] 計算用
- [ ] 条件式用
- [ ] 条件式グループ

*ループ*
- [ ] for
- [ ] while
- [ ] ループの中身
- [ ] 中断

### 一応定義
```text
MOLD = [STRING, INT, FLOAT, BOOL, MAP, LIST]  
VAR_DATA_TYPE = [INDENT, STRING, INT, FLOAT, TRUE, FALSE, MAP, LIST, NULL]  
CONTENT = [VarDef, Expr, WhileLoop, ForLoop, Return, CallFunc, ...]  
CALC_OPERATORS = [  
    PERCENT(%),  
    AST(*),  
    PLUS(+),  
    MINUS(-),  
    SLASH(/),   
]  
CONTROL_OPERATORS = [
    INCREMENT(++),  
    DECREMENT(--),  
    PLUSEq(+=),  
    MINUSEq(-=), 
]
COND_OPERATORS = [  
    EQUALEq(==),  
    QUESTEq(!=),  
    LT(<),  
    GT(>),  
    LTEq(<=),  
    GTEq(>=),  
    AND(&&),  
    OR(||),  
]  

LOGICAL_OPERATORS = [  
    AND(&&),  
    OR(||),  
]
```

#### RULES
- "+" : one or more  
- "*" : zero or more

## Variable
### VarDef(変数定義)

```text
Node{
    typ: VarDef
    childrenToken: [
        Token{ typ: IDENT }
    ]
    childrenNode : [
        ( VarData, CalcExpr )
    ]
}
```

### VarSubst(代入)
```text
Node{
    typ: Subst
    childrenToken: [
        Token{ typ: IDENT }
    ]
    childrenNode : [
        ( VarData, CalcExpr )
    ]
}
```

### VarData
```text
Node{
    typ: Data
    childrenToken: [
        Token{ typ: ( VAR_DATA_TYPE ) }
    ]
    childrenNode : []
}
```


## Contents(本文)
```text
Node{
    typ: Contents
    childrenToken: []
    childrenNode : [
        ( CONTENT )*
    ]
}
```

## Expr

### CalcExpr(数式)
```text
Node{
    typ: CalcExpr
    childrenToken: [
        ( CALC_OPERATORS )
        ( IDENT, INT, FLOAT )+
    ]
    childrenNode : []
}
```

### ControlExpr(値操作系の数式)
```text
Node{
    typ: CalcExpr
    childrenToken: [
        ( CONTROL_OPERATORS )
        ( IDENT, INT, FLOAT )+
    ]
    childrenNode : []
}
```

### CondExpr(条件式単体)
```text
Node{
    typ: CondExpr
    childrenToken: [
        ( COND_OPERATORS )
        ( VAR_DATA_TYPE )
        ( VAR_DATA_TYPE )
    ]
    childrenNode : []
}
```

### CondExprGroup(条件式グループ)
childrenTokenの個数は、CondExpr - 1、接続詞なので。
```text
Node{
    typ: CondExpr
    childrenToken: [
        ( LOGICAL_OPERATORS )*
    ]
    childrenNode : [
        ( CondExpr )+
    ]
}
```


## Function
### FuncDef(関数定義)

```text
Node{
    typ: FuncDef
    childrenToken: [
        Token{ typ: IDENT }
    ]
    childrenNode : [
        ( FuncFormalArgs )
        ( FuncRetType )
        ( FuncBody )
    ]
}
```

### FuncFormalArgs(仮引数)

```text
Node{ 
    typ: FuncFormalArgs
    childrenToken: []
    childrenNode : [
        ( FuncArg )*
    ]
}
```

### FuncArg(名前つき引数)
```text
Node{ 
    typ: FuncArg
    childrenToken: [
        Token{ typ: IDENT },
        Token{ typ: (MOLD) }
    ]
    childrenNode : []
}
```

### FuncRetTypes(戻り値)
```text
Node{ 
    typ: FuncRetTypes
    childrenToken: [
        Token{ typ: (MOLD) }*
    ]
    childrenNode : []
}
```

### FuncBody(関数本体)
```text
Node{ 
    typ: FuncBody
    childrenToken: []
    childrenNode : [
        ( Contents )
    ]
}
```

## Import
```text
Node{
    typ: Import
    childrenToken: [
        Token{ typ: STRING }
    ]
    childrenNode : []
}
```