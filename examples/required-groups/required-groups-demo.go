package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/internal/types"
)

/**
 * 必需组示例程序
 * 演示普通必需组和条件性必需组的使用
 * 通过子命令展示不同的使用场景
 */
func main() {
	// 配置全局根命令
	qflag.Root.SetDesc("必需组示例程序")
	qflag.Root.SetVersion("1.0.0")

	// 创建简单示例子命令
	simpleCmd := qflag.NewCmd("simple", "s", qflag.ContinueOnError)
	simpleCmd.SetDesc("简单必需组示例")

	// 创建普通必需组的标志
	serverHost := simpleCmd.String("server-host", "sh", "服务器主机", "")
	serverPort := simpleCmd.Uint("server-port", "sp", "服务器端口", 0)

	// 创建条件性必需组的标志
	dbHost := simpleCmd.String("db-host", "dh", "数据库主机", "")
	dbPort := simpleCmd.Uint("db-port", "dp", "数据库端口", 0)
	dbName := simpleCmd.String("db-name", "dn", "数据库名称", "")

	// 创建可选标志
	verbose := simpleCmd.Bool("verbose", "v", "详细输出", false)

	// 添加普通必需组 - 所有标志都必须设置
	err := simpleCmd.AddRequiredGroup("server", []string{"server-host", "server-port"}, false)
	if err != nil {
		fmt.Printf("添加普通必需组失败: %v\n", err)
		return
	}

	// 添加条件性必需组 - 如果使用其中一个则必须同时使用
	err = simpleCmd.AddRequiredGroup("database", []string{"db-host", "db-port", "db-name"}, true)
	if err != nil {
		fmt.Printf("添加条件性必需组失败: %v\n", err)
		return
	}

	// 设置简单示例的执行函数
	simpleCmd.SetRun(func(cmd types.Command) error {
		// 显示结果
		fmt.Printf("服务器: %s:%d\n", serverHost.Get(), serverPort.Get())

		if dbHost.IsSet() {
			fmt.Printf("数据库: %s:%d/%s\n", dbHost.Get(), dbPort.Get(), dbName.Get())
		} else {
			fmt.Println("数据库: 未配置")
		}

		if verbose.Get() {
			fmt.Println("详细模式: 已启用")
		}
		return nil
	})

	// 创建完整示例子命令
	fullCmd := qflag.NewCmd("full", "f", qflag.ContinueOnError)
	fullCmd.SetDesc("完整必需组示例")

	// 创建普通必需组的标志
	hostFlag := fullCmd.String("host", "h", "服务器主机地址", "")
	portFlag := fullCmd.Uint("port", "p", "服务器端口号", 0)

	// 创建条件性必需组的标志
	dbHostFlag := fullCmd.String("dbhost", "dh", "数据库主机地址", "")
	dbPortFlag := fullCmd.Uint("dbport", "dp", "数据库端口号", 0)
	dbUserFlag := fullCmd.String("dbuser", "du", "数据库用户名", "")
	_ = fullCmd.String("dbpass", "dw", "数据库密码", "") // 使用下划线忽略未使用的变量

	// 创建可选标志
	verboseFlag := fullCmd.Bool("verbose", "v", "详细输出", false)

	// 添加普通必需组 - 所有标志都必须设置
	err = fullCmd.AddRequiredGroup("server", []string{"host", "port"}, false)
	if err != nil {
		fmt.Printf("添加必需组失败: %v\n", err)
		return
	}

	// 添加条件性必需组 - 如果使用其中一个则必须同时使用
	err = fullCmd.AddRequiredGroup("database", []string{"dbhost", "dbport", "dbuser", "dbpass"}, true)
	if err != nil {
		fmt.Printf("添加条件性必需组失败: %v\n", err)
		return
	}

	// 设置完整示例的执行函数
	fullCmd.SetRun(func(cmd types.Command) error {
		// 使用参数
		fmt.Printf("服务器连接: %s:%d\n", hostFlag.Get(), portFlag.Get())

		if dbHostFlag.IsSet() {
			fmt.Printf("数据库连接: %s:%d, 用户: %s\n",
				dbHostFlag.Get(), dbPortFlag.Get(), dbUserFlag.Get())
		} else {
			fmt.Println("未使用数据库连接")
		}

		if verboseFlag.Get() {
			fmt.Println("详细模式已启用")
		}
		return nil
	})

	// 创建高级示例子命令
	advancedCmd := qflag.NewCmd("advanced", "a", qflag.ContinueOnError)
	advancedCmd.SetDesc("高级必需组示例 - 云服务部署工具")

	// 基础认证参数（普通必需组）
	accessKeyFlag := advancedCmd.String("access-key", "a", "访问密钥", "")
	secretKeyFlag := advancedCmd.String("secret-key", "s", "秘密密钥", "")

	// 基础部署参数（普通必需组）
	regionFlag := advancedCmd.String("region", "r", "部署区域", "")
	serviceFlag := advancedCmd.String("service", "svc", "服务名称", "")

	// 数据库连接参数（条件性必需组）
	advDbHostFlag := advancedCmd.String("db-host", "dh", "数据库主机", "")
	advDbPortFlag := advancedCmd.Uint("db-port", "dp", "数据库端口", 0)
	advDbNameFlag := advancedCmd.String("db-name", "dn", "数据库名称", "")
	advDbUserFlag := advancedCmd.String("db-user", "du", "数据库用户", "")
	advDbPassFlag := advancedCmd.String("db-pass", "dpw", "数据库密码", "")

	// Redis连接参数（条件性必需组）
	redisHostFlag := advancedCmd.String("redis-host", "rh", "Redis主机", "")
	redisPortFlag := advancedCmd.Uint("redis-port", "rp", "Redis端口", 0)
	redisPassFlag := advancedCmd.String("redis-pass", "rpw", "Redis密码", "")

	// 监控配置参数（条件性必需组）
	metricsHostFlag := advancedCmd.String("metrics-host", "mh", "监控主机", "")
	metricsPortFlag := advancedCmd.Uint("metrics-port", "mp", "监控端口", 0)
	metricsPathFlag := advancedCmd.String("metrics-path", "mpath", "监控路径", "")

	// 可选参数
	advVerboseFlag := advancedCmd.Bool("verbose", "v", "详细输出", false)
	dryRunFlag := advancedCmd.Bool("dry-run", "d", "试运行模式", false)

	// 添加普通必需组 - 认证参数
	err = advancedCmd.AddRequiredGroup("auth", []string{"access-key", "secret-key"}, false)
	if err != nil {
		fmt.Printf("添加认证必需组失败: %v\n", err)
		return
	}

	// 添加普通必需组 - 基础部署参数
	err = advancedCmd.AddRequiredGroup("deployment", []string{"region", "service"}, false)
	if err != nil {
		fmt.Printf("添加部署必需组失败: %v\n", err)
		return
	}

	// 添加条件性必需组 - 数据库连接
	err = advancedCmd.AddRequiredGroup("database", []string{"db-host", "db-port", "db-name", "db-user", "db-pass"}, true)
	if err != nil {
		fmt.Printf("添加数据库条件性必需组失败: %v\n", err)
		return
	}

	// 添加条件性必需组 - Redis连接
	err = advancedCmd.AddRequiredGroup("redis", []string{"redis-host", "redis-port", "redis-pass"}, true)
	if err != nil {
		fmt.Printf("添加Redis条件性必需组失败: %v\n", err)
		return
	}

	// 添加条件性必需组 - 监控配置
	err = advancedCmd.AddRequiredGroup("monitoring", []string{"metrics-host", "metrics-port", "metrics-path"}, true)
	if err != nil {
		fmt.Printf("添加监控条件性必需组失败: %v\n", err)
		return
	}

	// 设置高级示例的执行函数
	advancedCmd.SetRun(func(cmd types.Command) error {
		// 显示配置信息
		fmt.Println("=== 云服务部署配置 ===")
		fmt.Printf("访问密钥: %s\n", maskString(accessKeyFlag.Get()))
		fmt.Printf("秘密密钥: %s\n", maskString(secretKeyFlag.Get()))
		fmt.Printf("部署区域: %s\n", regionFlag.Get())
		fmt.Printf("服务名称: %s\n", serviceFlag.Get())

		// 显示数据库配置（如果设置了）
		if advDbHostFlag.IsSet() {
			fmt.Println("\n--- 数据库配置 ---")
			fmt.Printf("主机: %s:%d\n", advDbHostFlag.Get(), advDbPortFlag.Get())
			fmt.Printf("数据库: %s\n", advDbNameFlag.Get())
			fmt.Printf("用户: %s\n", advDbUserFlag.Get())
			fmt.Printf("密码: %s\n", maskString(advDbPassFlag.Get()))
		} else {
			fmt.Println("\n--- 数据库配置: 未启用 ---")
		}

		// 显示Redis配置（如果设置了）
		if redisHostFlag.IsSet() {
			fmt.Println("\n--- Redis配置 ---")
			fmt.Printf("主机: %s:%d\n", redisHostFlag.Get(), redisPortFlag.Get())
			fmt.Printf("密码: %s\n", maskString(redisPassFlag.Get()))
		} else {
			fmt.Println("\n--- Redis配置: 未启用 ---")
		}

		// 显示监控配置（如果设置了）
		if metricsHostFlag.IsSet() {
			fmt.Println("\n--- 监控配置 ---")
			fmt.Printf("主机: %s:%d%s\n", metricsHostFlag.Get(), metricsPortFlag.Get(), metricsPathFlag.Get())
		} else {
			fmt.Println("\n--- 监控配置: 未启用 ---")
		}

		// 显示可选参数
		if advVerboseFlag.Get() {
			fmt.Println("\n--- 其他选项 ---")
			fmt.Println("详细模式: 已启用")
		}

		if dryRunFlag.Get() {
			fmt.Println("试运行模式: 已启用")
		}

		fmt.Println("\n=== 部署准备完成 ===")
		return nil
	})

	// 添加子命令到全局根命令
	if err := qflag.AddSubCmds(simpleCmd, fullCmd, advancedCmd); err != nil {
		fmt.Printf("添加子命令失败: %v\n", err)
		return
	}

	// 解析并路由到子命令
	if err := qflag.ParseAndRoute(); err != nil {
		fmt.Printf("参数解析错误: %v\n", err)
		fmt.Println("\n使用示例:")
		fmt.Println("  简单示例: required-groups-demo simple --server-host localhost --server-port 8080")
		fmt.Println("  完整示例: required-groups-demo full --host localhost --port 8080")
		fmt.Println("  高级示例: required-groups-demo advanced --access-key AK123 --secret-key SK456 --region us-west-1 --service myapp")
		fmt.Println("\n使用 'required-groups-demo <subcommand> --help' 查看子命令的详细帮助")
		os.Exit(1)
	}
}

// maskString 隐藏敏感信息
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}
