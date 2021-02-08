<div align="center">
  <img src="gengine.png">
</div>

[![license](https://img.shields.io/badge/license-BSD-blue.svg)]()
[![Documentation](https://img.shields.io/badge/api-reference-blue.svg)](https://rencalo770.github.io/gengine_doc) 

# Gengine
- [English document](README.md)
- 使用交流QQ群号:1132683357 (解答所有的gengine问题)

## 基于golang的规则引擎
- **Gengine**是一款基于AST(Abstract Syntax Tree)和golang语言实现的规则引擎(动态化加载框架)。能够让你在golang这种静态语言上，在不停服务的情况下实现动态加载与配置规则。
- **代码结构松散，逻辑极其简单，但经过了必要且详尽的测试**
- Gengine所支持的规则，就是一门**DSL**(领域专用语言)
- Gengine是一款高性能框架动态加载，与golang无缝对接，支持多种执行模式、支持规则池等满足各种业务场景需求

## 设计思想
- 可以看这篇文章: https://xie.infoq.cn/article/40bfff1fbca1867991a1453ac

## 官方文档
- 因为gengine被设计成为极易使用的框架(仅有4个关键的API)，所以很多人在使用之前根本不看文档
- 我衷心的希望您在使用gengine之前,可以认真的看看文档.文档不仅能帮助你更好的使用gengine,还能帮助你拓宽你在业务上思路
- https://rencalo770.github.io/gengine_doc

## 支持的执行模式
 ![avatar](exe_model.jpg)

## 使用
- go mod 或者 go vendor
- 建议使用最新版本

## 问题联系
- 提issue