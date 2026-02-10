package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockFlag 基础模拟标志实现
type MockFlag struct {
	name       string
	short      string
	desc       string
	flagType   types.FlagType
	value      any
	strValue   string
	isSet      bool
	isRequired bool
	isHidden   bool
	envVar     string
	enumValues []string
}

// NewMockFlagBasic 创建基础模拟标志
func NewMockFlagBasic(name, desc string) *MockFlag {
	return &MockFlag{
		name:       name,
		desc:       desc,
		flagType:   types.FlagTypeString,
		value:      "",
		strValue:   "",
		isSet:      false,
		isRequired: false,
		isHidden:   false,
		envVar:     "",
		enumValues: []string{},
	}
}

// NewMockFlag 创建指定类型的模拟标志
func NewMockFlag(name, short, desc string, flagType types.FlagType, defaultValue any) *MockFlag {
	return &MockFlag{
		name:       name,
		short:      short,
		desc:       desc,
		flagType:   flagType,
		value:      defaultValue,
		strValue:   formatValue(defaultValue),
		isSet:      false,
		isRequired: false,
		isHidden:   false,
		envVar:     "",
		enumValues: []string{},
	}
}

// NewMockBoolFlag 创建布尔模拟标志
func NewMockBoolFlag(name, short, desc string, defaultValue bool) *MockFlag {
	return &MockFlag{
		name:       name,
		short:      short,
		desc:       desc,
		flagType:   types.FlagTypeBool,
		value:      defaultValue,
		strValue:   formatValue(defaultValue),
		isSet:      false,
		isRequired: false,
		isHidden:   false,
		envVar:     "",
		enumValues: []string{},
	}
}

// NewMockEnumFlag 创建枚举模拟标志
func NewMockEnumFlag(name, short, desc string, defaultValue string, allowedValues []string) *MockFlag {
	return &MockFlag{
		name:       name,
		short:      short,
		desc:       desc,
		flagType:   types.FlagTypeEnum,
		value:      defaultValue,
		strValue:   defaultValue,
		isSet:      false,
		isRequired: false,
		isHidden:   false,
		envVar:     "",
		enumValues: allowedValues,
	}
}

// 实现 Flag 接口
func (f *MockFlag) Name() string         { return f.name }
func (f *MockFlag) LongName() string     { return f.name }
func (f *MockFlag) ShortName() string    { return f.short }
func (f *MockFlag) Desc() string         { return f.desc }
func (f *MockFlag) Type() types.FlagType { return f.flagType }

func (f *MockFlag) Set(value string) error {
	f.strValue = value
	f.value = value
	f.isSet = true
	return nil
}

func (f *MockFlag) GetDef() any { return f.value }
func (f *MockFlag) GetStr() string {
	if f.isSet {
		return f.strValue
	}
	return formatValue(f.value)
}

func (f *MockFlag) IsSet() bool { return f.isSet }
func (f *MockFlag) Reset() {
	f.isSet = false
	f.strValue = formatValue(f.value)
}

func (f *MockFlag) String() string { return f.GetStr() }

func (f *MockFlag) BindEnv(env string) { f.envVar = env }
func (f *MockFlag) GetEnvVar() string  { return f.envVar }

func (f *MockFlag) EnumValues() []string { return f.enumValues }
func (f *MockFlag) Default() string      { return formatValue(f.value) }

func (f *MockFlag) IsRequired() bool { return f.isRequired }
func (f *MockFlag) IsHidden() bool   { return f.isHidden }

func (f *MockFlag) SetRequired(required bool) { f.isRequired = required }
func (f *MockFlag) SetHidden(hidden bool)     { f.isHidden = hidden }

// 辅助函数
func formatValue(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
