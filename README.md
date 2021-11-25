# Clock

***人生短短数十载 来都来了 凑活着过吧***

### 🔰 安装

使用 `go get` 或者直接下载预编译二进制文件 [clock/releases](https://github.com/chenjiandongx/clock/releases)

```shell
go get -u github.com/chenjiandongx/clock
```

### 🔑 配置读取

clock 默认会优先读取下列的环境变量

| Key | Desc | Default |
| --- | ---- | ------- |
| CLOCK_WHO | 称呼 | Coder |
| CLOCK_BIRTHDAY | 你来的日子 | 1996-04-12 |
| CLOCK_PASS_AWAY | 你走的日子 | 2086-04-12 |
| CLOCK_START_COLOR | 渐变起始颜色 | #F3C8ED |
| CLOCK_END_COLOR | 渐变终止颜色 | #B2F6EF |

新增 .life_clock 文件，以 yaml 格式配置，读取的优先级如下：

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

> 真正的勇士敢于直面赤裸裸的人生以及血淋淋的事实

![image](https://user-images.githubusercontent.com/19553554/143471486-0f9d2ab1-7756-492a-85de-b17ced0f3515.png)

### 👏 Contributors

- [NickChenyx](https://github.com/nickChenyx)

### 🔖 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
