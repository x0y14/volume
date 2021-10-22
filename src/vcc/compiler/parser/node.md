MOLD = [STRING, INT, FLOAT, BOOL, MAP, LIST]


*FuncDef*

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

*FuncFormalArgs*

```text
Node{ 
    typ: FuncFormalArgs
    childrenToken: []
    childrenNode : [
        ...
        ( FuncArg )
        ...
    ]
}
```

*FuncArg*
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

*FuncRetType*
```text
Node{ 
    typ: FuncRetType
    childrenToken: [
        ...
        Token{ typ: (MOLD) }
        ...
    ]
    childrenNode : []
}
```


<!--
Returnに統合、FuncBodyないに格納予定
*FuncRetType*
```text
Node{ 
    typ: FuncRetType
    childrenToken: [
        ...
        Token{ typ: (MOLD) }
        ...
    ]
    childrenNode : []
}
```
-->

*FuncBody*
```text
Node{ 
    typ: FuncBody
    childrenToken: []
    childrenNode : [
        ...
        todo
        ...
    ]
}
```