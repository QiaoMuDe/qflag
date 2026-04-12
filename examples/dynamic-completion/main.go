package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建根命令
	root := qflag.NewCmd("dynatest", "dt", qflag.ExitOnError)
	root.SetDesc("动态补全测试工具 - 用于测试 __complete 子命令和模糊匹配功能")
	root.SetCompletion(true)
	root.SetDynamicCompletion(true)

	// 添加全局标志
	root.String("config", "c", "配置文件路径", "")
	root.String("output", "o", "输出格式 (json|yaml|table)", "")
	root.Bool("verbose", "v", "启用详细输出", false)
	root.Bool("debug", "d", "启用调试模式", false)
	root.Enum("kind", "k", "资源类型", "service", []string{"service", "pod"})

	// 创建 service 子命令组
	serviceCmd := qflag.NewCmd("service", "svc", qflag.ExitOnError)
	serviceCmd.SetDesc("服务管理相关命令")
	serviceCmd.String("namespace", "n", "命名空间", "default")

	// service list 子命令
	serviceListCmd := qflag.NewCmd("list", "ls", qflag.ExitOnError)
	serviceListCmd.SetDesc("列出所有服务")
	serviceListCmd.String("filter", "f", "过滤条件", "")
	serviceListCmd.Bool("all", "a", "显示所有服务", false)
	serviceListCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("列出服务...")
		return nil
	})

	// service create 子命令
	serviceCreateCmd := qflag.NewCmd("create", "new", qflag.ExitOnError)
	serviceCreateCmd.SetDesc("创建新服务")
	serviceCreateCmd.String("name", "", "服务名称", "")
	serviceCreateCmd.String("type", "t", "服务类型", "clusterip")
	serviceCreateCmd.Int("port", "p", "服务端口号", 80)
	serviceCreateCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("创建服务...")
		return nil
	})

	// service delete 子命令
	serviceDeleteCmd := qflag.NewCmd("delete", "del", qflag.ExitOnError)
	serviceDeleteCmd.SetDesc("删除服务")
	serviceDeleteCmd.String("name", "", "要删除的服务名称", "")
	serviceDeleteCmd.Bool("force", "f", "强制删除", false)
	serviceDeleteCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("删除服务...")
		return nil
	})

	// service update 子命令
	serviceUpdateCmd := qflag.NewCmd("update", "up", qflag.ExitOnError)
	serviceUpdateCmd.SetDesc("更新服务")
	serviceUpdateCmd.String("name", "", "服务名称", "")
	serviceUpdateCmd.String("image", "i", "新镜像", "")
	serviceUpdateCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("更新服务...")
		return nil
	})

	// service logs 子命令
	serviceLogsCmd := qflag.NewCmd("logs", "log", qflag.ExitOnError)
	serviceLogsCmd.SetDesc("查看服务日志")
	serviceLogsCmd.String("name", "", "服务名称", "")
	serviceLogsCmd.Int("tail", "t", "显示最后N行", 100)
	serviceLogsCmd.Bool("follow", "f", "持续跟踪日志", false)
	serviceLogsCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("查看日志...")
		return nil
	})

	// 添加 service 子命令
	if err := serviceCmd.AddSubCmds(serviceListCmd, serviceCreateCmd, serviceDeleteCmd, serviceUpdateCmd, serviceLogsCmd); err != nil {
		fmt.Fprintf(os.Stderr, "添加 service 子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 deployment 子命令组
	deploymentCmd := qflag.NewCmd("deployment", "deploy", qflag.ExitOnError)
	deploymentCmd.SetDesc("部署管理相关命令")
	deploymentCmd.String("namespace", "n", "命名空间", "default")

	// deployment list 子命令
	deployListCmd := qflag.NewCmd("list", "ls", qflag.ExitOnError)
	deployListCmd.SetDesc("列出所有部署")
	deployListCmd.String("selector", "l", "标签选择器", "")
	deployListCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("列出部署...")
		return nil
	})

	// deployment create 子命令
	deployCreateCmd := qflag.NewCmd("create", "new", qflag.ExitOnError)
	deployCreateCmd.SetDesc("创建新部署")
	deployCreateCmd.String("name", "", "部署名称", "")
	deployCreateCmd.String("image", "i", "容器镜像", "")
	deployCreateCmd.Int("replicas", "r", "副本数量", 1)
	deployCreateCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("创建部署...")
		return nil
	})

	// deployment scale 子命令
	deployScaleCmd := qflag.NewCmd("scale", "", qflag.ExitOnError)
	deployScaleCmd.SetDesc("扩缩容部署")
	deployScaleCmd.String("name", "", "部署名称", "")
	deployScaleCmd.Int("replicas", "r", "目标副本数", 0)
	deployScaleCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("扩缩容部署...")
		return nil
	})

	// deployment rollback 子命令
	deployRollbackCmd := qflag.NewCmd("rollback", "rb", qflag.ExitOnError)
	deployRollbackCmd.SetDesc("回滚部署")
	deployRollbackCmd.String("name", "", "部署名称", "")
	deployRollbackCmd.Int("revision", "", "回滚到指定版本", 0)
	deployRollbackCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("回滚部署...")
		return nil
	})

	// 添加 deployment 子命令
	if err := deploymentCmd.AddSubCmds(deployListCmd, deployCreateCmd, deployScaleCmd, deployRollbackCmd); err != nil {
		fmt.Fprintf(os.Stderr, "添加 deployment 子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 config 子命令组
	configCmd := qflag.NewCmd("config", "cfg", qflag.ExitOnError)
	configCmd.SetDesc("配置管理相关命令")

	// config get 子命令
	configGetCmd := qflag.NewCmd("get", "", qflag.ExitOnError)
	configGetCmd.SetDesc("获取配置值")
	configGetCmd.String("key", "k", "配置键名", "")
	configGetCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("获取配置...")
		return nil
	})

	// config set 子命令
	configSetCmd := qflag.NewCmd("set", "", qflag.ExitOnError)
	configSetCmd.SetDesc("设置配置值")
	configSetCmd.String("key", "k", "配置键名", "")
	configSetCmd.String("value", "v", "配置值", "")
	configSetCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("设置配置...")
		return nil
	})

	// config list 子命令
	configListCmd := qflag.NewCmd("list", "ls", qflag.ExitOnError)
	configListCmd.SetDesc("列出所有配置")
	configListCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("列出配置...")
		return nil
	})

	// 添加 config 子命令
	if err := configCmd.AddSubCmds(configGetCmd, configSetCmd, configListCmd); err != nil {
		fmt.Fprintf(os.Stderr, "添加 config 子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 completion 子命令
	completionCmd := qflag.NewCmd("completion", "", qflag.ExitOnError)
	completionCmd.SetDesc("生成自动补全脚本")
	completionCmd.String("shell", "s", "Shell类型 (bash|pwsh|powershell)", "bash")
	completionCmd.SetRun(func(c qflag.Command) error {
		flag, exists := c.GetFlag("shell")
		if !exists {
			return fmt.Errorf("无法获取 shell 标志")
		}
		script, err := qflag.GenerateCompletion(root, flag.GetStr())
		if err != nil {
			return err
		}
		fmt.Println(script)
		return nil
	})

	// 创建更多子命令用于测试智能纠错
	// status 子命令
	statusCmd := qflag.NewCmd("status", "st", qflag.ExitOnError)
	statusCmd.SetDesc("查看系统状态")
	statusCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("系统状态: 正常运行")
		return nil
	})

	// version 子命令
	versionCmd := qflag.NewCmd("version", "v", qflag.ExitOnError)
	versionCmd.SetDesc("显示版本信息")
	versionCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("版本: v1.0.0")
		return nil
	})

	// info 子命令
	infoCmd := qflag.NewCmd("info", "i", qflag.ExitOnError)
	infoCmd.SetDesc("显示系统信息")
	infoCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("系统信息")
		return nil
	})

	// init 子命令
	initCmd := qflag.NewCmd("init", "", qflag.ExitOnError)
	initCmd.SetDesc("初始化配置")
	initCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("初始化完成")
		return nil
	})

	// validate 子命令
	validateCmd := qflag.NewCmd("validate", "val", qflag.ExitOnError)
	validateCmd.SetDesc("验证配置")
	validateCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("配置验证通过")
		return nil
	})

	// add 子命令
	addCmd := qflag.NewCmd("add", "a", qflag.ExitOnError)
	addCmd.SetDesc("添加资源")
	addCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("资源已添加")
		return nil
	})

	// remove 子命令
	removeCmd := qflag.NewCmd("remove", "rm", qflag.ExitOnError)
	removeCmd.SetDesc("移除资源")
	removeCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("资源已移除")
		return nil
	})

	// list 子命令（根命令级别）
	listCmd := qflag.NewCmd("list", "ls", qflag.ExitOnError)
	listCmd.SetDesc("列出所有资源")
	listCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("资源列表")
		return nil
	})

	// describe 子命令
	describeCmd := qflag.NewCmd("describe", "desc", qflag.ExitOnError)
	describeCmd.SetDesc("描述资源详情")
	describeCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("资源详情")
		return nil
	})

	// edit 子命令
	editCmd := qflag.NewCmd("edit", "e", qflag.ExitOnError)
	editCmd.SetDesc("编辑资源")
	editCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("进入编辑模式")
		return nil
	})

	// apply 子命令
	applyCmd := qflag.NewCmd("apply", "", qflag.ExitOnError)
	applyCmd.SetDesc("应用配置")
	applyCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("配置已应用")
		return nil
	})

	// delete 子命令（根命令级别）
	deleteCmd := qflag.NewCmd("delete", "del", qflag.ExitOnError)
	deleteCmd.SetDesc("删除资源")
	deleteCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("资源已删除")
		return nil
	})

	// watch 子命令
	watchCmd := qflag.NewCmd("watch", "w", qflag.ExitOnError)
	watchCmd.SetDesc("监视资源变化")
	watchCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("开始监视...")
		return nil
	})

	// exec 子命令
	execCmd := qflag.NewCmd("exec", "", qflag.ExitOnError)
	execCmd.SetDesc("执行命令")
	execCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("执行命令")
		return nil
	})

	// port-forward 子命令
	portForwardCmd := qflag.NewCmd("port-forward", "pf", qflag.ExitOnError)
	portForwardCmd.SetDesc("端口转发")
	portForwardCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("端口转发已建立")
		return nil
	})

	// proxy 子命令
	proxyCmd := qflag.NewCmd("proxy", "", qflag.ExitOnError)
	proxyCmd.SetDesc("代理服务")
	proxyCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("代理服务已启动")
		return nil
	})

	// cp 子命令
	cpCmd := qflag.NewCmd("cp", "", qflag.ExitOnError)
	cpCmd.SetDesc("复制文件")
	cpCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("文件复制完成")
		return nil
	})

	// auth 子命令
	authCmd := qflag.NewCmd("auth", "", qflag.ExitOnError)
	authCmd.SetDesc("认证管理")
	authCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("认证信息")
		return nil
	})

	// login 子命令
	loginCmd := qflag.NewCmd("login", "", qflag.ExitOnError)
	loginCmd.SetDesc("登录系统")
	loginCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("登录成功")
		return nil
	})

	// logout 子命令
	logoutCmd := qflag.NewCmd("logout", "", qflag.ExitOnError)
	logoutCmd.SetDesc("退出登录")
	logoutCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("已退出登录")
		return nil
	})

	// whoami 子命令
	whoamiCmd := qflag.NewCmd("whoami", "", qflag.ExitOnError)
	whoamiCmd.SetDesc("显示当前用户")
	whoamiCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("当前用户: admin")
		return nil
	})

	// plugin 子命令
	pluginCmd := qflag.NewCmd("plugin", "", qflag.ExitOnError)
	pluginCmd.SetDesc("插件管理")
	pluginCmd.SetRun(func(c qflag.Command) error {
		fmt.Println("插件管理")
		return nil
	})

	// 添加所有子命令到根命令
	if err := root.AddSubCmds(
		serviceCmd, deploymentCmd, configCmd, completionCmd,
		statusCmd, versionCmd, infoCmd, initCmd, validateCmd,
		addCmd, removeCmd, listCmd, describeCmd, editCmd,
		applyCmd, deleteCmd, watchCmd, execCmd, portForwardCmd,
		proxyCmd, cpCmd, authCmd, loginCmd, logoutCmd, whoamiCmd,
		pluginCmd,
	); err != nil {
		fmt.Fprintf(os.Stderr, "添加根命令子命令失败: %v\n", err)
		os.Exit(1)
	}

	// 设置根命令执行函数
	root.SetRun(func(c qflag.Command) error {
		fmt.Println("动态补全测试工具")
		fmt.Println()
		fmt.Println("使用方法:")
		fmt.Println("  dynatest service list")
		fmt.Println("  dynatest deployment create --name myapp --image nginx")
		fmt.Println("  dynatest config get --key database.url")
		fmt.Println()
		fmt.Println("生成补全脚本:")
		fmt.Println("  dynatest completion --shell bash")
		fmt.Println("  dynatest completion --shell pwsh")
		return nil
	})

	// 解析并执行
	if err := root.ParseAndRoute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
