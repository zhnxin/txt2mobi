# TXT2MOBI

[README_ZH_cn.md](./examples/README_ZH_cn.md)

`txt2mobi` is a simple tool for converting txt file to Mobi format for kindle. It's Based on golang package [mobi](https://github.com/766b/mobi/).

`Note`: All testing were done on Kindle App for windows and Kindle Paperwhite (10th Gen)
_PS: Any html element in text file would make effect._

## Build

```
go get -u github.com/zhnxin/txt2mobi
```

## Usage

1. Prepare the cover,thumbnail and txt file

2. Create a config file witch you can refer to  [example.toml](./examples/example.toml)

3. Run the command:
```sh
[linux]$./txt2mobi -f example.toml [- o output_file_name.mobi] [-p]
#for Windows
[windows]$ txt2mobi.exe -f example.toml [- o output_file_name.mobi] [-p]
```

## termimal parameter
require:
- `-config`: config file

`or`

require:

- `-f`: input file
- `-title`: book title
- `-author`: book author
- `-cover`: book cover and thumbnail
- `-chapter`: chater title regexp pattern

options:

- `-subchapter`: subchapter title regexp pattern
- `compress`: is to compress the result
- `encoding:`: use gb18030 as default


other options:

- `-o`: output file name
- `-p`: is to use '\<p\>\</p\>' to pack every paragrahes.
- `-escape`: to disable html escape

## config file

```toml
title="Example"
author="zhengxin"
file="example.txt"
cover="cover_example.jpg"
chapter="^Chapter\\.\\d+.*$"
subchapter='^Chapter\\.\\d+\-\\d+ .*$'
compress=false
//default cover size
//cover_width=860
//cover_hight=1200
```

_Note:_
- `title`: If the output_filename is not given, the `title` would be used as default.
- `chapter`: A regexp pattern to determin the chapter line title.
- `compress`: If it's true, the output would be compressed witch make the file smaller but take mush time to finish.
- `encoding`: Options. Use `utf-8` as default. Also support `gb18030` and `gbk`.

## Need GUI？

This tool is ready simple but modifis all need for me. Email me if you really desire it.

## About BLANK

The prefixed blank would be ignore because of html decoding.
