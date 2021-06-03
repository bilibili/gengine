<div align="center">
  <img src="gengine.png">
</div>

[![license](https://img.shields.io/badge/license-BSD-blue.svg)]()
[![Documentation](https://img.shields.io/badge/api-reference-blue.svg)](https://github.com/bilibili/gengine/wiki)

# Gengine
- [简体中文](README_zh.md)

## the rule engine based on golang
- this is a rule engine(or code dynamic load framework) named **Gengine** based on golang and AST, it can help you to load your code(rules) to run while you did not need to restart your application.
- Gengine's code structure is Modular design, logic is easy to understand, and necessary testing！
- it is also a high performance engine, support many execute-model and rules pool for business, it is easy to use in distribute system.

## github wiki
- English https://rencalo770.github.io/gengine_en/#/introduce
- 中文文档 https://github.com/bilibili/gengine/wiki

## supported the execute model of rules
 ![avatar](exe_model.jpg)


### env
`go1.12.x` (and later)

## use
- please use the newest version!
- go mod or go vendo, go mod:
```go
require github.com/bilibili/gengine v1.5.6
```

## Question Connection
- write issue
