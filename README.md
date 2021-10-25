# volume
a programming language  
Javaのように、バーチャルマシン上で動作します。

右辺、単項しかサポートしてない。

## VVM (Volume Virtual Machine)
VBinを読み込んで、実行するバーチャルマシン

## volc (volume compiler)
- vasm_gen : `*.vol` to `*.vol.s`
- vbin_gen : `*.vol.s` to `*.vol.b`

## Volume : `*.vol`

## VAsm (Volume Assembly) : `*.vol.s`
**INFO 世界で広く使用されている`Assembly`とは別物です**
`*.vol`から、変換されたもの。
中間表現。  
もしかしたら、複数の`.vol`ファイルを使用するプログラムを組んだときにオブジェクトファイルのように扱われるかもしれないが、今のところ`VBin`を生成する用途でしか使われない。  

## VBin (Volume Binary) : `*.vol.b`
**INFO 世界で広く使用されている`Binary`とは別物です**  
`*.vol.s`から変換されたもの。  
`VAsm`とほぼ同じ。  
`VVM`はこれを読み、実行する。  
`VVM`が読み込みやすいように作られている。


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