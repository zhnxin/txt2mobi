# TXT2Mobi

一个基于golang package [mobi](https://github.com/766b/mobi/)的简单命令行工具。能够将txt文本文件转换为Mobi格式

`注意`: 所有测试都是使用KindleForPC(windows)和Paperwhite (10th Gen)

_ps: 文档中的所有html元素都会生效_

## 编译


```
go get -u github.com/zhnxin/txt2mobi
```

## Usage

1. 准备好封面(.jpg),文本文档(.txt)
2. 根据[示例](./examples/example.toml)创建配置文件(.toml)
3. 运行
```sh
[linux]$./txt2mobi -f example.toml [- o output_file_name.mobi] [-p]
[windows]$ txt2mobi.exe -f example.toml [- o output_file_name.mobi] [-p]
```

## termimal parameter
require:
- `-config`: config file

or

require:
- `-f`: 文档
- `-title`: 书名
- `-author`: 作者
- `-cover`: 封面及缩略图
- `-chapter`: 一级标题的正则表达式
options:
- `-subchapter`: 二级标题的正则表达式，默认为空
- `compress`: 是否压缩，默认不使用
- `encoding`： 文件编码，默认为gb18030


other options:

- `-o`: 输出文件名
- `-p`: 是否使用 '\<p\>\</p\>' 装饰段落.
- `-escape`: 关闭html转义

## 配置文件
```toml
title="Example"
author="zhengxin"
file="example.txt"
cover="cover_example.jpg"
chapter="^第.*卷 .*$"
subchapter="^第\\d+章 .*$"
compress=false
```

_Note_:
- `title` : 如果可选的输出文件名没有提供，就会使用title作为默认文件名
- `chapter`： 正则表达式，用于识别章节
- `compress`： 被置为true时，会压缩生成结果，但回消耗更多的运行时间
- `encoding`： 可选参数，默认值为`utf-8`,也支持`gb18030`和`gbk`。

## 需要GUI?

这个命令行工具虽然简单(简陋？),但已经能够满足我的所有需求。如果你非常有需求，可以email联系我，让我有足够的￥动力￥来完成它。

## 关于开头空格

由于html内部解析，英文的段首空格将被忽略。
