# Clock

[![GoDoc](https://godoc.org/github.com/chenjiandongx/clock?status.svg)](https://godoc.org/github.com/chenjiandongx/clock)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenjiandongx/clock)](https://goreportcard.com/report/github.com/chenjiandongx/clock)
[![License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

***须知少时凌云志 曾许人间第一流***

![image](https://user-images.githubusercontent.com/19553554/143896547-c3bed94e-bccd-4624-9c98-580fb8517a5e.png)

### 🔰 安装

使用 `go get` 或者直接下载预编译二进制文件 [clock/releases](https://github.com/chenjiandongx/clock/releases)

```shell
go get -u github.com/chenjiandongx/clock
```

### 🔑 配置读取

clock 默认会优先读取下列的环境变量，日期的时间格式为 `yyyy-mm-dd`

| Key | Desc | Default |
| --- | ---- | ------- |
| CLOCK_WHO | 称呼 | Coder |
| CLOCK_BIRTHDAY | 你来的日子 | 1996-04-12 |
| CLOCK_PASS_AWAY | 你走的日子 | 2086-04-12 |
| CLOCK_START_COLOR | 渐变起始颜色 | #F3C8ED |
| CLOCK_END_COLOR | 渐变终止颜色 | #B2F6EF |

或者新增 .life_clock 文件，以 yaml 格式配置，读取的优先级如下：

* 0x00: 读取环境变量
* 0x01: 读取 clock 程序运行路径下的 .life_clock 文件
* 0x02: 读取 $HOME/.life_clock 文件
* 0x03: 读取 /etc/.life_clock 文件

文件默认配置如下：

```yaml
# .life_clock file

CLOCK_WHO: "Coder"
CLOCK_BIRTHDAY: "1996-04-12"
CLOCK_PASS_AWAY: "2086-04-12"
CLOCK_START_COLOR: "#F3C8ED"
CLOCK_END_COLOR: "#B2F6EF" 
```

### 🤔 人生阿！

> 真正的勇士敢于直面赤裸裸的人生和血淋淋的事实 -- dongdong

***短短这一生***

https://user-images.githubusercontent.com/19553554/143770322-01bbb073-4fc7-414d-bae6-01131cebb387.mov

### 👏 Contributors

- [NickChenyx](https://github.com/nickChenyx)

### 🔖 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
