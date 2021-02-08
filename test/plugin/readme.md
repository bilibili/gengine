# 基于go plugin的完全动态加载实现

## 测试用例使用方法
1.在当前文件夹中执行命令  go build -buildmode=plugin -o=plugin_M_m.so plugin_superman.go
2.然后产生一个plugin_M_m.so文件, 文件名中 M是在plugin_superman.go中最后一行导出的变量名; m是在gengine中使用的指针



