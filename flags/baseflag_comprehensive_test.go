package flags

import (
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/qerr"
)

// MockValidator 模拟验证器，用于测试验证功能
type MockValidator struct {
	shouldFail bool
	errorMsg   string
}

func (m *MockValidator) Validate(value any) error {
	if m.shouldFail {
		return qerr.NewValidationError(m.errorMsg)
	}
	return nil
}

// TestBaseFlag_Init 测试BaseFlag的初始化功能
func TestBaseFlag_Init(t *testing.T) {
	t.Run("正常初始化", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"
		err := flag.Init("longname", "l", "usage", &value)

		if err != nil {
			t.Errorf("正常初始化应该成功，但得到错误: %v", err)
		}

		if flag.LongName() != "longname" {
			t.Errorf("长名称应为 'longname'，实际为 '%s'", flag.LongName())
		}

		if flag.ShortName() != "l" {
			t.Errorf("短名称应为 'l'，实际为 '%s'", flag.ShortName())
		}

		if flag.Usage() != "usage" {
			t.Errorf("用法说明应为 'usage'，实际为 '%s'", flag.Usage())
		}
	})

	t.Run("重复初始化", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"

		// 第一次初始化
		err := flag.Init("longname", "l", "usage", &value)
		if err != nil {
			t.Fatalf("第一次初始化失败: %v", err)
		}

		// 第二次初始化应该失败
		err2 := flag.Init("longname2", "l2", "usage2", &value)
		if err2 == nil {
			t.Error("重复初始化应该返回错误")
		}
	})

	t.Run("空名称初始化", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"

		err := flag.Init("", "", "usage", &value)
		if err == nil {
			t.Error("长短名称都为空时应该返回错误")
		}
	})

	t.Run("空指针初始化", func(t *testing.T) {
		flag := &BaseFlag[string]{}

		err := flag.Init("longname", "l", "usage", nil)
		if err == nil {
			t.Error("值指针为nil时应该返回错误")
		}
	})
}

// TestBaseFlag_NameMethods 测试名称相关方法
func TestBaseFlag_NameMethods(t *testing.T) {
	t.Run("优先返回长名称", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"
		if err := flag.Init("longname", "l", "usage", &value); err != nil {
			t.Fatalf("初始化标志失败: %v", err)
		}

		if flag.Name() != "longname" {
			t.Errorf("Name()应返回长名称 'longname'，实际为 '%s'", flag.Name())
		}
	})

	t.Run("长名称为空时返回短名称", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"
		err := flag.Init("", "l", "usage", &value)
		if err != nil {
			t.Fatalf("初始化标志失败: %v", err)
		}

		if flag.Name() != "l" {
			t.Errorf("长名称为空时Name()应返回短名称 'l'，实际为 '%s'", flag.Name())
		}
	})
}

// TestBaseFlag_EnvironmentVariable 测试环境变量绑定功能
func TestBaseFlag_EnvironmentVariable(t *testing.T) {
	flag := &BaseFlag[string]{}
	value := "test"
	err := flag.Init("longname", "l", "usage", &value)
	if err != nil {
		t.Fatalf("初始化标志失败: %v", err)
	}

	// 测试绑定环境变量
	result := flag.BindEnv("TEST_ENV")
	if result != flag {
		t.Error("BindEnv应该返回自身以支持链式调用")
	}

	if flag.GetEnvVar() != "TEST_ENV" {
		t.Errorf("环境变量名应为 'TEST_ENV'，实际为 '%s'", flag.GetEnvVar())
	}
}

// TestBaseFlag_Validator 测试验证器功能
func TestBaseFlag_Validator(t *testing.T) {
	t.Run("验证器通过", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"
		err := flag.Init("longname", "l", "usage", &value)
		if err != nil {
			t.Fatalf("初始化标志失败: %v", err)
		}

		validator := &MockValidator{shouldFail: false}
		flag.SetValidator(validator)

		err = flag.Set("newvalue")
		if err != nil {
			t.Fatalf("设置标志失败: %v", err)
		}
	})

	t.Run("验证器失败", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"
		err := flag.Init("longname", "l", "usage", &value)
		if err != nil {
			t.Fatalf("初始化标志失败: %v", err)
		}

		validator := &MockValidator{shouldFail: true, errorMsg: "验证失败"}
		flag.SetValidator(validator)

		err = flag.Set("newvalue")
		if err == nil {
			t.Error("验证器失败时Set应该返回错误")
		}
	})
}

// TestBaseFlag_ConcurrentAccess 测试并发访问安全性
func TestBaseFlag_ConcurrentAccess(t *testing.T) {
	flag := &BaseFlag[int]{}
	value := 0
	err := flag.Init("longname", "l", "usage", &value)
	if err != nil {
		t.Fatalf("初始化标志失败: %v", err)
	}

	var wg sync.WaitGroup
	numGoroutines := 100

	// 并发设置值
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(val int) {
			defer wg.Done()
			err := flag.Set(val)
			if err != nil {
				t.Errorf("设置标志失败: %v", err)
			}
		}(i)
	}

	// 并发读取值
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			flag.Get()
			flag.IsSet()
			flag.GetPointer()
		}()
	}

	// 并发重置
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			flag.Reset()
		}()
	}

	wg.Wait()
	// 如果没有panic，说明并发访问是安全的
}

// TestBaseFlag_DefaultValueHandling 测试默认值处理
func TestBaseFlag_DefaultValueHandling(t *testing.T) {
	t.Run("获取默认值", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "default"
		err := flag.Init("longname", "l", "usage", &value)
		if err != nil {
			t.Fatalf("初始化标志失败: %v", err)
		}

		if flag.GetDefault() != "default" {
			t.Errorf("默认值应为 'default'，实际为 '%s'", flag.GetDefault())
		}

		if flag.GetDefaultAny() != "default" {
			t.Errorf("默认值(any)应为 'default'，实际为 '%v'", flag.GetDefaultAny())
		}
	})

	t.Run("未设置值时返回默认值", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "default"
		flag.Init("longname", "l", "usage", &value)

		if flag.Get() != "default" {
			t.Errorf("未设置值时Get()应返回默认值 'default'，实际为 '%s'", flag.Get())
		}

		if flag.IsSet() {
			t.Error("未设置值时IsSet()应返回false")
		}
	})
}

// TestBaseFlag_SetAndReset 测试设置和重置功能
func TestBaseFlag_SetAndReset(t *testing.T) {
	flag := &BaseFlag[string]{}
	value := "default"
	flag.Init("longname", "l", "usage", &value)

	// 设置新值
	err := flag.Set("newvalue")
	if err != nil {
		t.Fatalf("设置值失败: %v", err)
	}

	if flag.Get() != "newvalue" {
		t.Errorf("设置后的值应为 'newvalue'，实际为 '%s'", flag.Get())
	}

	if !flag.IsSet() {
		t.Error("设置值后IsSet()应返回true")
	}

	// 重置值
	flag.Reset()

	if flag.Get() != "default" {
		t.Errorf("重置后的值应为默认值 'default'，实际为 '%s'", flag.Get())
	}

	if flag.IsSet() {
		t.Error("重置后IsSet()应返回false")
	}
}

// TestBaseFlag_PointerAccess 测试指针访问功能
func TestBaseFlag_PointerAccess(t *testing.T) {
	t.Run("未设置值时指针指向初始值", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "default"
		flag.Init("longname", "l", "usage", &value)

		ptr := flag.GetPointer()
		if ptr == nil {
			t.Fatal("GetPointer()不应返回nil")
		}

		// 未设置值时，指针应指向初始值
		if *ptr != "default" {
			t.Errorf("未设置值时指针应指向初始值 'default'，实际为 '%s'", *ptr)
		}

		// 验证IsSet状态
		if flag.IsSet() {
			t.Error("未设置值时IsSet()应返回false")
		}
	})

	t.Run("设置值后指针有效", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "default"
		flag.Init("longname", "l", "usage", &value)

		flag.Set("newvalue")
		ptr := flag.GetPointer()

		if ptr == nil {
			t.Fatal("设置值后GetPointer()不应返回nil")
		}

		if *ptr != "newvalue" {
			t.Errorf("指针指向的值应为 'newvalue'，实际为 '%s'", *ptr)
		}

		// 通过指针修改值
		*ptr = "modified"
		if flag.Get() != "modified" {
			t.Errorf("通过指针修改后的值应为 'modified'，实际为 '%s'", flag.Get())
		}
	})
}

// TestBaseFlag_StringRepresentation 测试字符串表示
func TestBaseFlag_StringRepresentation(t *testing.T) {
	t.Run("字符串类型", func(t *testing.T) {
		flag := &BaseFlag[string]{}
		value := "test"
		flag.Init("longname", "l", "usage", &value)

		err := flag.Set("hello")
		if err != nil {
			t.Fatalf("设置标志失败: %v", err)
		}
		if flag.String() != "hello" {
			t.Errorf("字符串表示应为 'hello'，实际为 '%s'", flag.String())
		}
	})

	t.Run("整数类型", func(t *testing.T) {
		flag := &BaseFlag[int]{}
		value := 0
		flag.Init("longname", "l", "usage", &value)

		err := flag.Set(42)
		if err != nil {
			t.Fatalf("设置标志失败: %v", err)
		}
		if flag.String() != "42" {
			t.Errorf("字符串表示应为 '42'，实际为 '%s'", flag.String())
		}
	})
}

// TestBaseFlag_TypeMethod 测试Type方法的默认实现
func TestBaseFlag_TypeMethod(t *testing.T) {
	flag := &BaseFlag[string]{}
	value := "test"
	flag.Init("longname", "l", "usage", &value)

	if flag.Type() != 0 {
		t.Errorf("BaseFlag的Type()默认实现应返回0，实际返回 %d", flag.Type())
	}
}
