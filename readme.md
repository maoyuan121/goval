goval
=====

[![Build Status](https://travis-ci.org/maja42/goval.svg?branch=master)](https://travis-ci.org/maja42/goval)
[![Go Report Card](https://goreportcard.com/badge/github.com/maja42/goval)](https://goreportcard.com/report/github.com/maja42/goval)
[![Coverage Status](https://coveralls.io/repos/github/maja42/goval/badge.svg?branch=master)](https://coveralls.io/github/maja42/goval?branch=master)
[![GoDoc](https://godoc.org/github.com/maja42/goval?status.svg)](https://godoc.org/github.com/maja42/goval)

这个库允许程序计算任意算术/字符串/逻辑表达式。
支持访问变量和调用自定义函数。

这个项目被认为是稳定的，并且已经在生产系统中使用。
对于任何问题、反馈和错误报告，请使用 bug-reports。

开源协议 MIT license。

# Demo

一个对表达式求值的 CLI 小演示可以在 example 文件夹中找到:

```
go get -u github.com/maja42/goval
cd $GOPATH/src/github.com/maja42/goval/
go run example/main.go
```

![Demo](goval.gif)


# Usage

小演示：

```go
eval := goval.NewEvaluator()
result, err := eval.Evaluate(`42 > 21`, nil, nil) // Returns <true, nil>
```

访问变量：
```go
eval := goval.NewEvaluator()
variables := map[string]interface{}{
    "uploaded": 146,
    "total":  400,
}
result, err := eval.Evaluate(`uploaded * 100 / total`, variables, nil)  // Returns <36, nil>
```

调用函数：
```go
// Implementing strlen()
eval := goval.NewEvaluator()
variables := map[string]interface{}{
    "os":   runtime.GOOS,
    "arch": runtime.GOARCH,
}

functions := make(map[string]goval.ExpressionFunction)
functions["strlen"] = func(args ...interface{}) (interface{}, error) {
    str := args[0].(string)
    return len(str), nil
}

result, err := eval.Evaluate(`strlen(arch[:2]) + strlen("text")`, variables, functions) // Returns <6, nil>
```

自定义函数允许扩展的任意特性，如 regex 匹配:
```go
// Implementing regular expressions (error handling omitted)
functions := make(map[string]goval.ExpressionFunction)
functions["matches"] = func(args ...interface{}) (interface{}, error) {
    str := args[0].(string)
    exp := args[1].(string)
    reg := regexp.MustCompile(exp)
    return reg.MatchString(str), nil
}

eval.Evaluate(`matches("text", "[a-z]+")`, nil, functions)  // Returns <true, nil>
eval.Evaluate(`matches("1234", "[a-z]+")`, nil, functions)  // Returns <false, nil>
```

# 文档

## 类型

这个库完全支持以下类型:`nil`, `bool`, `int`, `float64`, `string`, `[]interface{}` (=arrays) and `map[string]interface{}` (=objects). 

在表达式中，`int` 和 `float64` 都具有 `number` 类型，并且是完全透明的。\
如果需要，数值将在 `int` 和 `float64` 之间自动转换，只要不丢失精度。

数组和对象是无类型的。它们可以存储任何其他值(“混合数组”)。

支持 struct 以保持功能的清晰和可管理。
它们会引入太多的边缘情况和松散的结果，因此超出了范围。

## 变量

可以直接访问自定义变量。
变量是只读的，不能从表达式中修改。

实例：

```
var
var.field
var[0]
var["field"]
var[anotherVar]

var["fie" + "ld"].field[42 - var2][0]
```

## Functions

可以从表达式中调用自定义函数。

实例：

```
rand()
floor(42)
min(4, 3, 12, max(1, 3, 3))
len("te" + "xt")
```

## Literals

任何字面量都可以在表达式中定义。
字符串字面值可以放在双引号 `"` 或反勾号 \`中。
十六进制字面值以前缀 `0x` 开始。

实例：

```
nil
true
false
3
3.2
"Hello, 世界!\n"
"te\"xt"
`te"xt`
[0, 1, 2]
[]
[0, ["text", false], 4.2]
{}
{"a": 1, "b": {c: 3}}
{"key" + 42: "value"}
{"k" + "e" + "y": "value"}

0xA                 // 10
0x0A                // 10
0xFF                // 255 
0xFFFFFFFF          // 32bit appl.: -1  64bit appl.: 4294967295
0xFFFFFFFFFFFFFFFF  // 64bit appl.: -1  32bit appl.: error
```

It is possible to access elements of array and object literals:

Examples:

```
[1, 2, 3][1]                // 2
[1, [2, 3, 42][1][2]        // 42

{"a": 1}.a                  // 1
{"a": {"b": 42}}.a.b        // 42
{"a": {"b": 42}}["a"]["b"]  // 42
```

## 优先级

严格遵循操作符优先级 [C/C++ rules](http://en.cppreference.com/w/cpp/language/operator_precedence)。

括号 `()` 用于控制优先级。

例子：

```
1 + 2 * 3    // 7
(1 + 2) * 3  // 9
```

## 操作符

### 算术

#### 算术 `+` `-` `*` `/`

如果两边都是整数，则结果值也是整数。
否则，结果将是一个浮点数。

例子：

```
3 + 4               // 7
2 + 2 * 3           // 8
2 * 3 + 2.5         // 8.5
12 - 7 - 5          // 0
24 / 10             // 2
24.0 / 10           // 2.4
```

#### Modulo `%`

如果两边都是整数，则结果值也是整数。
否则，结果将是一个浮点数。

例子：

```
4 % 3       // 1
144 % 85    // -55
5.5 % 2     // 1.5
10 % 3.5    // 3.0
```

#### 负数 `-` (unary minus)

例子：

```
-4       // -4
5 + -4   // 1
-5 - -4  // -1
1 + --1  // syntax error
-(4+3)   // -7
-varName
```


### Concatenation

#### 字符串连接 `+`

如果 `+` 操作符的左边或右边有一个是 `string`，则执行字符串连接。
支持字符串、数字、布尔值和nil。

例子:

```
"text" + 42     // "text42"
"text" + 4.2    // "text4.2"
42 + "text"     // "42text"
"text" + nil    // "textnil"
"text" + true   // "texttrue"
```

#### 数组连接 `+`

如果 `+` 操作符的两边都是数组，则将它们连接起来

例子：

```
[0, 1] + [2, 3]          // [0, 1, 2, 3]
[0] + [1] + [[2]] + []   // [0, 1, [2]]
```

#### 对象连接 `+`

如果 `+` 操作符的两边都是对象，则它们的字段将组合成一个新对象。
如果两个对象包含相同的键，则右边对象的值将覆盖左边对象的值。

例子:

```
{"a": 1} + {"b": 2} + {"c": 3}         // {"a": 1, "b": 2, "c": 3}
{"a": 1, "b": 2} + {"b": 3, "c": 4}    // {"a": 1, "b": 3, "c": 4}
{"b": 3, "c": 4} + {"a": 1, "b": 2}    // {"a": 1, "b": 2, "c": 4}
```

### Logic

#### Equals `==`, NotEquals `!=`

在两个操作数之间进行深度比较。
比较 `int` 和 `float64` 时，该整数将被转换为浮点数。

#### Comparisons `<`, `>`, `<=`, `>=`

比较两个数字。如果运算符的一边是整数，另一边是浮点数，整数值将被转换。对于在此过程中四舍五入的非常大的数字，这可能会导致意想不到的结果。

例子:

```
3 <-4        // false
45 > 3.4     // false
-4 <= -1     // true
3.5 >= 3.5   // true
```

#### And `&&`, Or `||`

例子:

```
true && true             // true
false || false           // false
true || false && false   // true
false && false || true   // true
```


#### Not `!`

Inverts the boolean on the right.

例子：

```
!true       // false
!false      // true
!!true      // true
!varName
```


### Ternary `? :`

如果表达式解析为 `true`，则运算符解析为左操作数。＼
如果表达式解析为 `false`，则运算符解析为右操作数。

例子：

```
true  ? 1 : 2                         // 1
false ? 1 : 2                         // 2
	
2 < 5  ? "a" : 1.5                    // "a"
9 > 12 ? "a" : [42]                   // [42]

false ? (true ? 1:2) : (true ? 3:4)   // 3
```

请注意，所有操作数都已解析(没有短路)。

在下面的例子中，两个函数都被调用了(`func2` 的返回值被简单地忽略了):

```
true ? func1() : func2()
```

### Bit Manipulation

#### Logical Or `|`, Logical And `&`, Logical XOr `^`

If one side of the operator is a floating point number, the number is cast to an integer if possible. 
If decimal places would be lost during that process, it is considered a type error.
The resulting number is always an integer.

Examples:

```
8 | 2          // 10
9 | 5          // 13
8 | 2.0        // 10
8 | 2.1        // type error

13 & 10        // 8
10 & 15.0 & 2  // 2

13 ^ 10        // 7
10 ^ 15 ^ 1    // 4
```

#### Bitwise Not `~`

If performed on a floating point number, the number is cast to an integer if possible. 
If decimal places would be lost during that process, it is considered a type error.
The resulting number is always an integer.

The results can differ between 32bit and 64bit architectures.

Examples:

```
~-1                   // 0
(~0xA55A) & 0xFFFF    // 0x5AA5
(~0x5AA5) & 0xFFFF    // 0xA55A

~0xFFFFFFFF           // 64bit appl.: 0xFFFFFFFF 00000000; 32bit appl.: 0x00
~0xFFFFFFFF FFFFFFFF  // 64bit appl.: 0x00; 32bit: error
```

#### Bit-Shift `<<`, `>>`

If one side of the operator is a floating point number, the number is cast to an integer if possible. 
If decimal places would be lost during that process, it is considered a type error.
The resulting number is always an integer.

When shifting to the right, sign-extension is performed.
The results can differ between 32bit and 64bit architectures.

Examples:

```
1 << 0    // 1
1 << 1    // 2
1 << 2    // 4
8 << -1   // 4
8 >> -1   // 16

1 << 31   // 0x00000000 80000000   64bit appl.: 2147483648; 32bit appl.: -2147483648
1 << 32   // 0x00000001 00000000   32bit appl.: 0 (overflow)

1 << 63   // 0x80000000 00000000   32bit appl.: 0 (overflow); 64bit appl.: -9223372036854775808
1 << 64   // 0x00000000 00000000   0 (overflow)

0x80000000 00000000 >> 63     // 0xFFFFFFFF FFFFFFFF   64bit: -1 (sign extension); 32bit: error (cannot parse number literal)
0x80000000 >> 31              // 64bit: 0x00000000 0000001; 32bit: 0xFFFFFFFF (-1, sign extension)
```

### More

#### Array contains `in`

如果数组包含特定元素，则返回 true 或 false。

例子:

```
"txt" in [nil, "hello", "txt", 42]   // true
true  in [nil, "hello", "txt", 42]   // false
nil   in [nil, "hello", "txt", 42]   // true
42.0  in [nil, "hello", "txt", 42]   // true
2         in [1, [2, 3], 4]          // false
[2, 3]    in [1, [2, 3], 4]          // true
[2, 3, 4] in [1, [2, 3], 4]          // false
```

#### Substrings `[a:b]`

Slices a string and returns the given substring.
Strings are indexed byte-wise. Multi-byte characters need to be treated carefully.

The start-index indicates the first byte to be present in the substring.\
The end-index indicates the last byte NOT to be present in the substring.\
Hence, valid indices are in the range `[0, len(str)]`.

Examples:

```
"abcdefg"[:]    // "abcdefg"
"abcdefg"[1:]   // "bcdefg"
"abcdefg"[:6]   // "abcdef"
"abcdefg"[2:5]  // "cde"
"abcdefg"[3:4]  // "d"

// The characters 世 and 界 both require 3 bytes:
"Hello, 世界"[7:13]    // "世界"
"Hello, 世界"[7:10]    // "世"
"Hello, 世界"[10:13]   // "界"
```


#### Array Slicing `[a:b]`

Slices an array and returns the given subarray.

The start-index indicates the first element to be present in the subarray.\
The end-index indicates the last element NOT to be present in the subarray.\
Hence, valid indices are in the range `[0, len(arr)]`.

Examples:

```
// Assuming `arr := [0, 1, 2, 3, 4, 5, 6]`:
arr[:]    // [0, 1, 2, 3, 4, 5, 6]
arr[1:]   // [1, 2, 3, 4, 5, 6]
arr[:6]   // [0, 1, 2, 3, 4, 5]
arr[2:5]  // [2, 3, 4]
arr[3:4]  // [3]
```

# Alternative Libraries

If you are looking for a generic evaluation library, 
you can also take a look at [Knetic/govaluate](https://github.com/Knetic/govaluate).
I used that library myself, but due to a several shortcomings I decided to create goval. 
The main differences are:

- More intuitive syntax
- No intermediate AST - evaluation and parsing happens in a single step  
- Better type support:  
    - Full support for arrays and objects.
    - Opaque differentiation between `int` and `float64`. \
      The underlying type is automatically converted as long as no precision is lost.
    - Type-aware bit-operations (they only work with `int`-numbers).
    - No support for dates (strings are just strings, they don't have a special meaning, even if they look like dates).\
      Support for dates and structs *could* be added if needed.
- More operators:      
    - Accessing variables (maps) via `.` and `[]` syntax
    - Support for array- and object concatenation.
    - Slicing and substrings
- Hex-Literals (useful as soon as bit-operations are involved).
- Array literals with `[]` as well as object literals with `{}`
- Useful error messages.
- Highly optimized parser code by using go/scanner and goyacc. \
    This leads to vastly reduced code size (and therefore little bug potential)
    and creates super-fast code.
- High test coverage (including lots of special cases).\
  Also tested on 32 and 64bit architectures, where some (documented) operations like a bitwise-not can behave differently depending on the size of `int`. 


