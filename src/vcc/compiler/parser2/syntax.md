##Kind
- DataTypes
  - D_STRING
  - D_INT
  - D_FLOAT
  - D_BOOL
  - D_NULL
- ExprOperators (全部用意しちゃう:-)
  - ExprItem
- Import ライブラリインポート
- FuncDef 関数定義 
  - FuncParams 関数の仮引数
    - FuncParamsItem
      - FuncParamsItemData
  - FuncReturnValues
  - FuncBody = Scripts(ifとかと使いまわせるように)
- Scripts プログラム１行、１行の集合体
- Expression
    - ConditionalExpr 条件式
- If
  - ConditionalExpr
  - Body = Scripts
- Elseif
  - ConditionalExpr
  - Body = Scripts
- Else
  - Body = Scripts
- For
  - ConditionalExpr
  - Body = Scripts
- While
  - ConditionalExpr
  - Body = Scripts
- CallFunc 呼び出し
  - CallFuncArgs
- VarDef 変数
  - VarRHS
    - VarRHSItem

### template
```text
Node {
    kind : ,
    tok  : ,
    nods : [],
}
```

###Import
```text
Node {
    kind : Import,
    tok  : Token.STRING,
    nods : [],
}
```

##Function
###FuncDef
```text
Node {
    kind : FuncDef,
    tok  : Token.IDENT,
    nods : [],
}
```

###FuncParams(入れ物)
```text
Node {
    kind : FuncParams,
    tok  : ,
    nods : [FuncParamsItem, ...],
}
```

###FuncParamsItem
```text
Node {
    kind : FuncParamsItem,
    tok  : Token.IDENT,
    nods : [FuncParamsItemData],
}
```

###FuncParamsItemData
```text
Node {
    kind : (D_STRING | D_INT | D_FLOAT | D_BOOL | D_NULL),
    tok  : Token.Any,
    nods : [],
}
```

###FuncReturnValues(入れ物)
```text
Node {
    kind : FuncReturnValues,
    tok  : ,
    nods : [(D_STRING | D_INT | D_FLOAT | D_BOOL | D_NULL), ...],
}
```

###FuncBody = Scripts

##Scripts
###Scripts
```text
Node {
    kind : Scripts,
    tok  : ,
    nods : [Script, ...],
}
```

###Script = (VarDef, ...)

## Expression
### ExprItem
```text
Node {
    kind : ExprItem,
    tok  : 任意のトークン,
    nods : [],
}
```
### ConditionalExpr
真偽を求める式
```text
Node {
    kind : ConditionalExpr,
    tok  : ,
    nods : [ExprItem, ...],
}
```


##Conditional Branch
###If
```text
Node {
    kind : If,
    tok  : ,
    nods : [ConditionalExpr, Scripts],
}
```

###Elseif
```text
Node {
    kind : Elseif,
    tok  : ,
    nods : [ConditionalExpr, Scripts],
}
```

###Else
```text
Node {
    kind : Else,
    tok  : ,
    nods : [Scripts],
}
```

##Loop
###For
```text
Node {
    kind : For,
    tok  : ,
    nods : [ConditionalExpr, Scripts],
}
```

###While
```text
Node {
    kind : While,
    tok  : ,
    nods : [ConditionalExpr, Scripts],
}
```

##CallFunc
###CallFunc(呼び出し)
```text
Node {
    kind : CallFunc,
    tok  : ,
    nods : [CallFuncArgs],
}
```
###CallFuncArgs(入れ物)
```text
Node {
    kind : CallFuncArgs,
    tok  : ,
    nods : [CallFuncItem, ...],
}
```
###CallFuncItem
```text
Node {
    kind : CallFuncItem,
    tok  : Token.Any,
    nods : [],
}
```

##Variable
### VarDef
```text
Node {
    kind : VarDef,
    tok  : Token.IDENT,
    nods : [VarRHS],
}
```

###VarRHS
```text
Node {
    kind : VarRHS,
    tok  : ,
    nods : [VarItem, ...],
}
```

###VarRHSItem
```text
Node {
    kind : VarRHSItem,
    tok  : Token.Any,
    nods : [],
}
```
