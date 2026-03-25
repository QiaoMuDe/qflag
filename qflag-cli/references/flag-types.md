# QFlag 标志类型详解

## 基础类型

### StringFlag

```go
flag := cmd.String("name", "n", "描述", "默认值")

// 使用
value := flag.Get()           // 获取值
str := flag.GetStr()          // 获取字符串表示
isSet := flag.IsSet()         // 检查是否被设置
```

### BoolFlag

```go
flag := cmd.Bool("verbose", "v", "详细输出", false)

// 使用
value := flag.Get()
```

### IntFlag / Int64Flag

```go
flag := cmd.Int("count", "c", "数量", 0)
flag64 := cmd.Int64("size", "s", "大小", 0)

// 使用
value := flag.Get()
```

### UintFlag / Uint8 / Uint16 / Uint32 / Uint64

```go
flag := cmd.Uint("port", "p", "端口", 8080)
flag8 := cmd.Uint8("level", "l", "级别", 1)
flag16 := cmd.Uint16("code", "c", "代码", 0)
flag32 := cmd.Uint32("id", "i", "ID", 0)
flag64 := cmd.Uint64("total", "t", "总数", 0)
```

### Float64Flag

```go
flag := cmd.Float64("rate", "r", "比率", 0.0)
```

## 特殊类型

### EnumFlag

限制为预定义值集合。

```go
flag := cmd.Enum("format", "f", "输出格式", "json", []string{"json", "xml", "yaml"})

// 使用
format := flag.Get()
```

## 时间类型

### DurationFlag

支持时间单位：ns, us, ms, s, m, h

```go
flag := cmd.Duration("timeout", "t", "超时时间", time.Second*30)

// 使用
duration := flag.Get()  // time.Duration 类型
```

### TimeFlag

支持多种时间格式。

```go
flag := cmd.Time("start", "s", "开始时间", time.Now())

// 使用
t := flag.Get()  // time.Time 类型
```

## 大小类型

### SizeFlag

支持存储单位：B, KB, MB, GB, TB, PB

```go
flag := cmd.Size("limit", "l", "限制大小", 1024)

// 使用
size := flag.Get()  // int64 类型，单位为字节
```

## 集合类型

### StringSliceFlag

```go
flag := cmd.StringSlice("tags", "t", "标签列表", []string{"default"})

// 使用
tags := flag.Get()  // []string 类型
```

### IntSliceFlag / Int64SliceFlag

```go
flag := cmd.IntSlice("ports", "p", "端口列表", []int{8080, 8081})
flag64 := cmd.Int64Slice("ids", "i", "ID列表", []int64{})

// 使用
ports := flag.Get()  // []int 类型
```

### MapFlag

```go
flag := cmd.Map("env", "e", "环境变量", map[string]string{"KEY": "VALUE"})

// 使用
env := flag.Get()  // map[string]string 类型
```

## 标志通用方法

所有标志类型都支持以下方法：

```go
// 获取值
value := flag.Get()

// 获取字符串表示
str := flag.GetStr()

// 获取默认值
defaultValue := flag.GetDef()

// 检查是否被设置
isSet := flag.IsSet()

// 获取环境变量名
envVar := flag.GetEnvVar()

// 设置验证器
flag.SetValidator(validator)

// 重置为默认值
flag.Reset()

// 获取标志信息
name := flag.Name()           // 长名称
short := flag.ShortName()     // 短名称
desc := flag.Desc()           // 描述
flagType := flag.Type()       // 标志类型
```
