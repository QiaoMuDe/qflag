#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
禁用标志解析功能自动化测试脚本

本脚本用于自动化测试 qflag 的 DisableFlagParsing 功能，
验证各种场景下禁用标志解析的行为是否符合预期。

测试场景包括：
1. exec 子命令 - 透传参数给容器
2. ssh 子命令 - SSH 包装器
3. normal 子命令 - 正常解析标志（对比）
4. wrapper 子命令 - 方法链配置

使用方法：
    python test_disable_flag_parsing.py

环境要求：
    - Python 3.6+
    - Go 环境已配置
    - 位于 examples/disable-flag-parsing 目录
"""

import subprocess
import sys
import os
from typing import List, Tuple, Optional
from dataclasses import dataclass
from enum import Enum


class TestResult(Enum):
    """测试结果枚举"""
    PASS = "通过"
    FAIL = "失败"
    SKIP = "跳过"


@dataclass
class TestCase:
    """测试用例数据类"""
    name: str                    # 测试名称
    description: str             # 测试描述
    args: List[str]              # 命令参数
    expected_disable: bool       # 期望的禁用标志解析状态
    expected_in_output: List[str]  # 期望输出中包含的字符串
    not_expected_in_output: List[str] = None  # 不期望输出中包含的字符串


class Colors:
    """终端颜色代码"""
    GREEN = ''
    RED = ''
    YELLOW = ''
    BLUE = ''
    RESET = ''
    BOLD = ''
    
    @classmethod
    def init_colors(cls):
        """初始化颜色（在非Windows或支持颜色的终端）"""
        if sys.platform != 'win32' or 'ANSICON' in os.environ:
            cls.GREEN = '\033[92m'
            cls.RED = '\033[91m'
            cls.YELLOW = '\033[93m'
            cls.BLUE = '\033[94m'
            cls.RESET = '\033[0m'
            cls.BOLD = '\033[1m'


class DisableFlagParsingTester:
    """禁用标志解析功能测试器"""
    
    def __init__(self):
        self.test_cases: List[TestCase] = []
        self.passed = 0
        self.failed = 0
        self.skipped = 0
        self.script_dir = os.path.dirname(os.path.abspath(__file__))
        
    def add_test(self, test_case: TestCase):
        """添加测试用例"""
        self.test_cases.append(test_case)
        
    def run_command(self, args: List[str]) -> Tuple[int, str, str]:
        """
        执行命令并返回结果
        
        参数:
            args: 命令参数列表
            
        返回:
            (返回码, 标准输出, 标准错误)
        """
        cmd = ['go', 'run', 'main.go'] + args
        try:
            result = subprocess.run(
                cmd,
                cwd=self.script_dir,
                capture_output=True,
                text=True,
                encoding='utf-8',
                timeout=30
            )
            return result.returncode, result.stdout, result.stderr
        except subprocess.TimeoutExpired:
            return -1, "", "命令执行超时"
        except Exception as e:
            return -1, "", str(e)
            
    def check_output(self, output: str, expected: List[str], not_expected: Optional[List[str]] = None) -> Tuple[bool, List[str]]:
        """
        检查输出是否包含期望的字符串
        
        参数:
            output: 命令输出
            expected: 期望包含的字符串列表
            not_expected: 不期望包含的字符串列表
            
        返回:
            (是否通过, 错误信息列表)
        """
        errors = []
        
        # 检查期望包含的字符串
        for exp in expected:
            if exp not in output:
                errors.append(f"期望包含 '{exp}'，但未找到")
                
        # 检查不期望包含的字符串
        if not_expected:
            for not_exp in not_expected:
                if not_exp in output:
                    errors.append(f"不期望包含 '{not_exp}'，但找到了")
                    
        return len(errors) == 0, errors
        
    def print_header(self, text: str):
        """打印标题"""
        print(f"\n{Colors.BOLD}{Colors.BLUE}{'='*60}{Colors.RESET}")
        print(f"{Colors.BOLD}{Colors.BLUE}{text.center(60)}{Colors.RESET}")
        print(f"{Colors.BOLD}{Colors.BLUE}{'='*60}{Colors.RESET}\n")
        
    def print_test_start(self, test_name: str, description: str):
        """打印测试开始信息"""
        print(f"{Colors.BOLD}测试: {test_name}{Colors.RESET}")
        print(f"描述: {description}")
        
    def print_test_result(self, result: TestResult, details: str = ""):
        """打印测试结果"""
        if result == TestResult.PASS:
            print(f"{Colors.GREEN}[PASS] 通过{Colors.RESET}")
        elif result == TestResult.FAIL:
            print(f"{Colors.RED}[FAIL] 失败{Colors.RESET}")
        else:
            print(f"{Colors.YELLOW}[SKIP] 跳过{Colors.RESET}")
            
        if details:
            print(f"  {details}")
        print()
        
    def run_test(self, test_case: TestCase) -> TestResult:
        """执行单个测试用例"""
        self.print_test_start(test_case.name, test_case.description)
        print(f"命令: go run main.go {' '.join(test_case.args)}")
        
        # 执行命令
        returncode, stdout, stderr = self.run_command(test_case.args)
        output = stdout + stderr
        
        # 检查命令是否成功执行
        if returncode != 0:
            self.print_test_result(TestResult.FAIL, f"命令执行失败，返回码: {returncode}\n错误: {stderr}")
            return TestResult.FAIL
            
        # 检查禁用标志解析状态 (Go 输出小写 true/false)
        expected_status = f"禁用标志解析: {str(test_case.expected_disable).lower()}"
        if expected_status not in output:
            self.print_test_result(TestResult.FAIL, f"未找到 '{expected_status}'\n实际输出:\n{output}")
            return TestResult.FAIL
            
        # 检查期望输出
        passed, errors = self.check_output(
            output, 
            test_case.expected_in_output,
            test_case.not_expected_in_output
        )
        
        if not passed:
            self.print_test_result(TestResult.FAIL, "\n".join(errors))
            return TestResult.FAIL
            
        self.print_test_result(TestResult.PASS)
        return TestResult.PASS
        
    def run_all_tests(self):
        """执行所有测试"""
        self.print_header("禁用标志解析功能自动化测试")
        
        print(f"测试用例总数: {len(self.test_cases)}")
        print(f"脚本目录: {self.script_dir}")
        print()
        
        for i, test_case in enumerate(self.test_cases, 1):
            print(f"[{i}/{len(self.test_cases)}] ", end="")
            result = self.run_test(test_case)
            
            if result == TestResult.PASS:
                self.passed += 1
            elif result == TestResult.FAIL:
                self.failed += 1
            else:
                self.skipped += 1
                
        self.print_summary()
        
    def print_summary(self):
        """打印测试摘要"""
        self.print_header("测试摘要")
        
        total = len(self.test_cases)
        pass_rate = (self.passed / total * 100) if total > 0 else 0
        
        print(f"总测试数: {total}")
        print(f"{Colors.GREEN}通过: {self.passed}{Colors.RESET}")
        print(f"{Colors.RED}失败: {self.failed}{Colors.RESET}")
        print(f"{Colors.YELLOW}跳过: {self.skipped}{Colors.RESET}")
        print(f"通过率: {pass_rate:.1f}%")
        
        if self.failed == 0:
            print(f"\n{Colors.GREEN}{Colors.BOLD}所有测试通过！[OK]{Colors.RESET}")
        else:
            print(f"\n{Colors.RED}{Colors.BOLD}存在失败的测试，请检查！[ERROR]{Colors.RESET}")
            
        print()


def create_test_cases() -> List[TestCase]:
    """创建测试用例列表"""
    return [
        # ========== exec 子命令测试 ==========
        TestCase(
            name="exec - 基本透传",
            description="测试 exec 子命令透传基本参数",
            args=["exec", "mypod", "--", "ls", "-la"],
            expected_disable=True,
            expected_in_output=[
                "=== exec 子命令 ===",
                "禁用标志解析: true",
                "原始参数: [mypod -- ls -la]",
                "Pod 和选项: [mypod]",
                "执行命令: [ls -la]"
            ]
        ),
        TestCase(
            name="exec - 无分隔符",
            description="测试 exec 子命令无 -- 分隔符的情况",
            args=["exec", "mypod", "ls", "-la"],
            expected_disable=True,
            expected_in_output=[
                "禁用标志解析: true",
                "原始参数: [mypod ls -la]",
            ]
        ),
        TestCase(
            name="exec - 空参数",
            description="测试 exec 子命令无参数时的帮助信息",
            args=["exec"],
            expected_disable=True,
            expected_in_output=[
                "禁用标志解析: true",
                "用法: demo exec <pod名> [--] <命令> [参数...]"
            ]
        ),
        TestCase(
            name="exec - 复杂参数",
            description="测试 exec 子命令处理复杂参数",
            args=["exec", "mypod", "-c", "mycontainer", "--", "ps", "aux"],
            expected_disable=True,
            expected_in_output=[
                "禁用标志解析: true",
                "原始参数: [mypod -c mycontainer -- ps aux]"
            ]
        ),
        
        # ========== ssh 子命令测试 ==========
        TestCase(
            name="ssh - 基本透传",
            description="测试 ssh 子命令透传参数",
            args=["ssh", "user@192.168.1.1", "--", "ls", "-la"],
            expected_disable=True,
            expected_in_output=[
                "=== ssh 子命令 ===",
                "禁用标志解析: true",
                "原始参数: [user@192.168.1.1 -- ls -la]",
                "将执行: ssh user@192.168.1.1 -- ls -la"
            ]
        ),
        TestCase(
            name="ssh - SSH选项",
            description="测试 ssh 子命令处理 SSH 选项",
            args=["ssh", "user@host", "-p", "2222", "--", "ps", "aux"],
            expected_disable=True,
            expected_in_output=[
                "禁用标志解析: true",
                "原始参数: [user@host -p 2222 -- ps aux]"
            ]
        ),
        
        # ========== normal 子命令测试（对比）==========
        TestCase(
            name="normal - 解析标志",
            description="测试 normal 子命令正常解析标志",
            args=["normal", "--verbose", "--output=test.txt", "arg1", "arg2"],
            expected_disable=False,
            expected_in_output=[
                "=== normal 子命令 ===",
                "禁用标志解析: false",
                "verbose: true",
                "output: test.txt",
                "位置参数: [arg1 arg2]"
            ]
        ),
        TestCase(
            name="normal - 短标志",
            description="测试 normal 子命令解析短标志",
            args=["normal", "-v", "-o", "out.log"],
            expected_disable=False,
            expected_in_output=[
                "禁用标志解析: false",
                "verbose: true"
            ]
        ),
        TestCase(
            name="normal - 仅位置参数",
            description="测试 normal 子命令无标志的情况",
            args=["normal", "arg1", "arg2"],
            expected_disable=False,
            expected_in_output=[
                "禁用标志解析: false",
                "verbose: false",
                "位置参数: [arg1 arg2]"
            ]
        ),
        
        # ========== wrapper 子命令测试 ==========
        TestCase(
            name="wrapper - 方法链配置",
            description="测试 wrapper 子命令通过方法链禁用标志解析",
            args=["wrapper", "--any-flag", "value", "--another"],
            expected_disable=True,
            expected_in_output=[
                "=== wrapper 子命令 (通过方法链配置) ===",
                "禁用标志解析: true",
                "原始参数: [--any-flag value --another]",
                "此命令通过 SetDisableFlagParsing(true) 禁用标志解析"
            ]
        ),
        TestCase(
            name="wrapper - 所有参数作为位置参数",
            description="验证 wrapper 子命令所有参数都作为位置参数",
            args=["wrapper", "--flag1", "val1", "-f", "val2", "posarg"],
            expected_disable=True,
            expected_in_output=[
                "禁用标志解析: true",
                "原始参数: [--flag1 val1 -f val2 posarg]"
            ]
        ),
    ]


def main():
    """主函数"""
    # 检查是否在正确的目录
    script_dir = os.path.dirname(os.path.abspath(__file__))
    main_go = os.path.join(script_dir, "main.go")
    
    if not os.path.exists(main_go):
        print(f"{Colors.RED}错误: 未找到 main.go 文件{Colors.RESET}")
        print(f"请确保在 examples/disable-flag-parsing 目录下运行此脚本")
        sys.exit(1)
        
    # 初始化颜色
    Colors.init_colors()
    
    # 检查 Go 环境
    try:
        result = subprocess.run(['go', 'version'], capture_output=True, text=True)
        if result.returncode != 0:
            raise Exception("Go 命令执行失败")
        print(f"{Colors.GREEN}Go 环境: {result.stdout.strip()}{Colors.RESET}\n")
    except Exception as e:
        print(f"{Colors.RED}错误: 无法检测到 Go 环境 - {e}{Colors.RESET}")
        sys.exit(1)
        
    # 创建测试器并运行测试
    tester = DisableFlagParsingTester()
    
    # 添加所有测试用例
    for test_case in create_test_cases():
        tester.add_test(test_case)
        
    # 运行测试
    tester.run_all_tests()
    
    # 返回退出码
    sys.exit(0 if tester.failed == 0 else 1)


if __name__ == "__main__":
    main()
