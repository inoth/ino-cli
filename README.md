# ino-cli

```shell
ubuntu@ubuntu% ino-cli help 
inoth的工具箱，看心情添加功能。
author： inoth

Usage:
  ino-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gen         读取数据库结构生成实体
  help        Help about any command
  init        初始化项目手脚架
  rest        解析swagger文档，生成http文件
  version     show version detail.

Flags:
  -h, --help     help for ino-cli
  -t, --toggle   Help message for toggle

Use "ino-cli [command] --help" for more information about a command.
```
## 功能点扩充计划
* 手脚架功能，整体生成后考虑后续新增常用中间件，实现 add 命令添加 redis、mongodb、mysql 之类的，注册后直接开始使用
