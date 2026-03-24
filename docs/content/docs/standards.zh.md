---
title: '文档章节标准'
linkTitle: 'standards'
description: 'ArcGo 各子包文档必须包含的章节清单'
weight: 2
draft: false
---

## 适用范围

本标准适用于 ArcGo 的核心子包文档：

- `authx`、`clientx`、`collectionx`、`configx`、`dbx`、`dix`、`eventx`、`httpx`、`kvx`、`logx`、`observabilityx`、`sqltmplx`

`examples/*` 目录仅作为可运行示例代码来源，不作为独立“子包文档体系”。

## 每个子包必须具备的章节

每个子包入口页（`_index.md`）至少包含以下章节（建议按顺序）：

1. **Overview / 包定位**
   - 说明该包是什么
   - 说明它与相邻子包的边界
2. **Install / Import**
   - `go get` 路径
   - 如有可选子模块，也要写明
3. **Quick Start**
   - 最小可运行示例（含完整 import）
4. **Core Capabilities**
   - 核心能力清单
5. **Key API Surface**
   - 高频类型 / 函数
   - 推荐默认路径 API
6. **Configuration and Options**
   - 选项、预设、默认行为、扩展点
7. **Error and Behavior Model**
   - 错误分类与行为契约
8. **Integration Guide**
   - 与其他 ArcGo 子包如何组合
9. **Testing and Production Notes**
   - 测试建议、benchmark 命令、上线注意点
10. **Examples**
   - 示例命令与路径
   - 示例仅为支撑代码，不单独作为子包体系

## 内容规则

- 代码片段必须包含完整 import。
- 优先写可操作示例，不堆概念描述。
- 子包文档中不再保留 roadmap / iteration plan 内容。
- 有多语言页时，`*.md`、`*.en.md`、`*.zh.md` 结构应尽量对齐。

