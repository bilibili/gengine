package test

import "fmt"

type UnitPrintMock struct {
	debugInfo map[int]string
	count     int
}

// NewUnitPrintMock 用于初始化单测的的print mock
func NewUnitPrintMock() *UnitPrintMock {
	return &UnitPrintMock{debugInfo: map[int]string{}, count: 0}
}

// Debug 用于mock脚本中的println
func (imp *UnitPrintMock) Debug(a ...interface{}) {
	imp.debugInfo[imp.count] = fmt.Sprint(a...)
	fmt.Println(a...)
	imp.count++
}

// Debugf 用于mock脚本中的println
func (imp *UnitPrintMock) Debugf(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	imp.debugInfo[imp.count] = str
	fmt.Println(str)
	imp.count++
}

// Get 获取打印的信息
func (imp *UnitPrintMock) Get(key int) string {
	return imp.debugInfo[key]
}
