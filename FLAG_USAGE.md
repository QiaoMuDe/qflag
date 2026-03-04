# 标志使用语法指南

本文档详细说明了QFlag库中标志的使用语法，包括基础语法和高级标志的用法。

## 基础语法

QFlag基于Go标准库的`flag`包，支持以下标志语法：

### 1. 短标志名（单字符）

```bash
# 使用空格分隔
-f value

# 使用等号分隔
-f=value
```

### 2. 长标志名（多字符）

```bash
# 使用空格分隔
--flag value

# 使用等号分隔
--flag=value
```

### 3. 布尔标志

布尔标志不需要指定值，存在即为true：

```bash
# 短布尔标志
-f

# 长布尔标志
--flag

# 可以使用等号显式设置（可选）
--flag=true
--flag=false
```

### 4. 组合使用

```bash
# 同时使用多个标志
-f value1 --flag2 value2 -g value3
```

## 高级标志语法

除了基础标志外，QFlag还提供了多种高级标志类型。

### 1. 时间标志（Duration）

时间标志接受多种时间格式：

```bash
# 直接指定数字（默认为纳秒）
--timeout 1000000000  # 1秒

# 使用时间单位
--timeout 1s           # 1秒
--timeout 100ms         # 100毫秒
--timeout 5m           # 5分钟
--timeout 2h           # 2小时

# 复合时间
--timeout 1h30m        # 1小时30分钟
--timeout 2d5h30m     # 2天5小时30分钟
```

支持的时间单位：
- `ns` - 纳秒
- `us` - 微秒
- `ms` - 毫秒
- `s` - 秒
- `m` - 分钟
- `h` - 小时
- `d` - 天

### 2. 大小标志（Size）

大小标志接受多种大小格式：

```bash
# 直接指定数字（默认为字节）
--size 1024           # 1024字节

# 使用大小单位
--size 1KB            # 1千字节
--size 10MB           # 10兆字节
--size 2GB            # 2吉字节

# 使用小数
--size 1.5MB         # 1.5兆字节
```

支持的大小单位：
- `B` - 字节
- `KB` - 千字节
- `MB` - 兆字节
- `GB` - 吉字节
- `TB` - 太字节

### 3. 切片标志（Slice）

切片标志可以接受多个值：

```bash
# 多次使用同一标志
--file file1.txt --file file2.txt --file file3.txt

# 使用逗号分隔
--file file1.txt,file2.txt,file3.txt

# 混合使用
--file file1.txt --file file2.txt,file3.txt --file file4.txt
```

### 4. 映射标志（Map）

映射标志接受键值对：

```bash
# 使用等号分隔键值
--param key1=value1 --param key2=value2

# 使用多次指定
--param key1=value1,key2=value2,key3=value3
```

### 5. 枚举标志（Enum）

枚举标志限制为预定义的值：

```bash
# 使用预定义的值
--mode debug          # 正确
--mode production     # 正确
--mode invalid         # 错误，会报错
```

## 标志定义示例

```go
// 基础标志
name := cmd.String("name", "n", "用户名", "default")
age := cmd.Int("age", "a", "年龄", 18)
enabled := cmd.Bool("enabled", "e", "是否启用", false)

// 高级标志
timeout := cmd.Duration("timeout", "t", "超时时间", time.Second*30)
size := cmd.Size("limit", "l", "大小限制", 10*1024*1024) // 10MB
files := cmd.StringSlice("files", "f", "文件列表", []string{})
config := cmd.Map("config", "c", "配置参数", map[string]string{})
mode := cmd.Enum("mode", "m", "运行模式", []string{"debug", "release"}, "release")
```

## 使用建议

1. **一致性** - 在整个项目中保持一致的命名和风格
2. **默认值** - 为标志提供合理的默认值
3. **帮助文本** - 编写清晰、简洁的帮助文本
4. **验证** - 使用验证器确保输入值的有效性
5. **分组** - 使用互斥组和必需组组织相关标志

## 常见错误

1. **未识别的标志** - 检查标志名称是否正确
2. **缺少参数** - 非布尔标志需要提供值
3. **格式错误** - 检查时间、大小等特殊格式的正确性
4. **类型不匹配** - 确保提供的值可以转换为标志类型

## 高级技巧

1. **环境变量** - 可以绑定环境变量作为默认值
2. **子命令** - 使用子命令组织复杂的功能
3. **条件解析** - 根据某些标志的存在决定其他标志的行为
4. **自定义验证** - 添加自定义验证逻辑确保数据完整性