世界的に有名なGo言語の専門家として、ご提供いただいたPDFファイルの内容に基づき、Markdown形式のチートシートを生成します。

***

# Go言語 チートシート

このチートシートは、Go言語の基本的な構文と概念をまとめたものです。

## 変数 (Variables)

### 基本的な宣言
```go
[cite_start]// 型を指定して宣言 [cite: 1]
var a string

[cite_start]// 宣言と同時に初期化 [cite: 1]
var a string = "something"

[cite_start]// 型推論による宣言 [cite: 1]
var a = "something"

[cite_start]// 短縮形式での宣言（関数内でのみ使用可能） [cite: 1]
a := "something"

[cite_start]// 複数の変数を同時に宣言 [cite: 1]
a, b := "one", "two"

[cite_start]// varブロックでまとめて宣言 [cite: 1]
var (
    a string
    b int
)
```

### Map
```go
[cite_start]// mapの宣言 [cite: 1]
var a map[string]string

[cite_start]// make関数でmapを初期化 [cite: 1]
a := make(map[string]string)

[cite_start]// mapリテラルで初期化 [cite: 1]
a := map[string]string{"k": "v"}
```

### Array
```go
[cite_start]// 配列の宣言と初期化 [cite: 1]
var d [2]int = [2]int{1, 2}
```

### Slice
```go
[cite_start]// sliceの宣言 [cite: 1]
var a []string

[cite_start]// make関数でsliceを初期化（長さ5） [cite: 1]
a := make([]string, 5)

[cite_start]// make関数でsliceを初期化（長さ5、容量10） [cite: 1]
a := make([]string, 5, 10)

[cite_start]// sliceリテラルで初期化 [cite: 1]
a := []int{1, 2, 3, 4, 5, 6, 7}
```
* [cite_start]インデックス`0`から`2`の手前まで（`[1 2]`） [cite: 1]
    ```go
    fmt.Printf("%v", a[0:2])
    ```
* [cite_start]インデックス`5`から最後まで（`[6 7]`） [cite: 1]
    ```go
    fmt.Printf("%v", a[5:len(a)])
    ```

### Struct
```go
[cite_start]// 構造体の型を定義 [cite: 1]
type myStruct struct {
    c string
}

[cite_start]// 構造体の変数を宣言 [cite: 1]
var a myStruct

[cite_start]// 構造体を初期化 [cite: 1]
a := myStruct{c: "a string"}

[cite_start]// フィールドに値を代入 [cite: 1]
a.c = "a string"
```

## 型 (Types)
```
[cite_start]bool [cite: 1]
[cite_start]string [cite: 1]
[cite_start]byte (uint8のエイリアス) [cite: 1]
[cite_start]rune (int32のエイリアス) [cite: 1]
[cite_start]float32 float64 [cite: 1]
[cite_start]int int8 int16 int32 int64 [cite: 1]
[cite_start]uint uint8 uint16 uint32 uint64 uintptr [cite: 1]
[cite_start]complex64 complex128 [cite: 1]
```

## 演算子 (Operators)
* [cite_start]**算術演算子**: `a := 5 + 5` [cite: 1]
* [cite_start]**代入演算子**: `a += 5` (`a = a + 5`と同じ) [cite: 1]
* [cite_start]**比較演算子**: `b := a >= 5` (`true`または`false`を返す) [cite: 1]
* [cite_start]**インクリメント**: `b++` (bを1増やす) [cite: 1]
* [cite_start]**デクリメント**: `b--` (bを1減らす) [cite: 1]

## 条件 (Conditions)

### 比較・論理演算子
* [cite_start]**比較**: `< > <= >= == !=` [cite: 1]
* **論理**: `|| [cite_start]&& !` [cite: 2]

### if文
```go
if a >= 10 && b >= 10 {
    // ...
} else if !(a < 5) {
    // ...
} else {
    // ...
}
```

### switch文
```go
switch a {
case 5:
    // ...
[cite_start]default: // オプショナル [cite: 1]
    // ...
}
```

## ループ (Loops)

### for文
```go
[cite_start]// 基本的なforループ [cite: 3]
for i := 0; i < 5; i++ {
    // ...
}

[cite_start]// 無限ループ [cite: 3]
for {
    // ...
}

[cite_start]// break: ループを抜ける [cite: 3]
if a == 5 {
    break
}

[cite_start]// continue: 次のイテレーションへ進む [cite: 3]
if a == 5 {
    continue
}
```

### rangeキーワード
```go
[cite_start]// 配列、スライス、マップの要素をループ [cite: 4]
for k, v := range a {
    [cite_start]// 配列/スライスの場合、kはインデックス(int)、vは値 [cite: 4]
    [cite_start]// マップの場合、kはキー、vは値 [cite: 4]
}
```

## 関数 (Functions)
```go
[cite_start]// 基本的な関数 [cite: 1]
func myFunc(b int, c int64) bool {
    return true
}

[cite_start]// 引数の型が同じ場合は省略可能 [cite: 1]
func myFunc(b, c int) { /* ... */ }

[cite_start]// 複数の戻り値 [cite: 1]
func returnTwo() (bool, bool) { /* ... */ }

[cite_start]// 関数を変数に代入（無名関数） [cite: 1]
a := func(a int) int { return a + 1 }
fmt.Printf("%v", a(1))
```

## メソッド (Methods)
```go
type myStruct struct {
    c string
}

[cite_start]// 値レシーバ（呼び出し元の値は変更されない） [cite: 4]
func (m myStruct) method() {
    m.c = "newstr"
}

[cite_start]// ポインタレシーバ（呼び出し元の値が変更される） [cite: 4]
func (m *myStruct) method() {
    m.c = "newstr"
}
```

## ポインタ (Pointers)
```go
a := "123"

[cite_start]// ポインタを取得（aのアドレスをbに代入） [cite: 4]
b := &a

[cite_start]fmt.Println(b) // 例: 0xc00001c030（アドレス） [cite: 4]
[cite_start]fmt.Println(*b) // "123"（ポインタが指す値） [cite: 4]

[cite_start]// ポインタを通じて元の変数の値を変更 [cite: 4]
*b = "456"
```

## インターフェース (Interfaces)
```go
[cite_start]// メソッドのシグネチャを定義 [cite: 4]
type myIface interface {
    method()
    method2(a int) (bool, error)
}
```

## その他 (Various)

### 定数
```go
[cite_start]// 値が固定の定数を宣言 [cite: 4]
const a int = 5
```

### ゴルーチン
```go
[cite_start]// hello()関数を非同期に実行 [cite: 4]
go hello()
```