<div align="center">
  <img src="gengine.png">
</div>

[![license](https://img.shields.io/badge/license-BSD-blue.svg)]()
[![Documentation](https://img.shields.io/badge/api-reference-blue.svg)](https://rencalo770.github.io/gengine_en) 

# Gengine
- [简体中文](README_zh.md) 使用交流Q群1132683357

## the rule engine based on golang 
- this is a rule engine(or code dynamic load framework) named **Gengine** based on golang and AST, it can help you to load your code(rules) to run while you did not need to restart your application.  
- Gengine's code structure is Modular design, logic is easy to understand, and necessary testing！
- it is also a high performance engine, support many execute-model and rules pool for business, it is easy to use in distribute system. 

## English Doc
- because gengine is designed to be extremely easy to use(only 4 APIs), so there are a lot programmers to use gengine without reading document!
- I sincerely hope you can read the document before you use gengine! The doc can greatly help you to use it well
- https://rencalo770.github.io/gengine_en

## supported the execute model of rules
 ![avatar](exe_model.jpg)

## use 
- go mod or go vendor 
- please use the newest version! 
```go
module your_module
go 1.14
require github.com/bilibili/gengine v1.5.0
```


## Question Connection
- write issue