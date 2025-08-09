package utils // import "gitee.com/MM-Q/qflag/utils"

utils 工具包 Package utils 通用工具函数集合 本文件提供了qflag包使用的各种通用工具函数，包括字符串处理、类型转换、
文件操作等辅助功能，为其他模块提供基础支持。

FUNCTIONS

func GetExecutablePath() string
    GetExecutablePath 获取程序的绝对安装路径 如果无法通过 os.Executable 获取路径,则使用 os.Args[0] 作为替代

    返回值:
      - 程序的绝对路径字符串

