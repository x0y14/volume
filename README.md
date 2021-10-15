# volume
a programming language

Javaのように、Virtual Machine上で動作します。

VMは擬似機械語を使用します。

- 擬似アセンブリ(`*.vol.s`): `*.vol`を解析し、アセンブリのような言語に変換したもの。  
機械語に変換しやすくなっている。  
一見、本物のアセンブリのように見えるが、Golangの言語機能に頼り切った命令が多く存在する。

- 擬似機械語(`出力なし`): 擬似アセンブリで、使用されていた変数等々を、数値に変換し、読み込みやすくしたもの。   
ただし、本当の機械語のように0, 1で記述されているわけではない。  
あくまで、擬似機械語。  
ミスリードを誘うような命名だが、これ以上最適なものを考えられなかったためこれを使用する。



cp src dst  
add src dst  
sub src dst  


#### string idea
- len "text" dst  
  dst = 4

- split "text"  
push "t"  
push "e"  
push "x"  
push "t"  

- join "text" dst  
  dst = dst + "text"  