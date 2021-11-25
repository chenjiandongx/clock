# Clock

***人生短短数十载 来都来了 凑活着过吧***

### 🔰 安装

使用 `go get` 或者直接下载预编译二进制文件 [clock/release](https://github.com/chenjiandongx/clock/releases)

```shell
go get -u github.com/chenjiandongx/clock
```

### 🔑 环境变量

clock 默认会优先读取下列的环境变量

| Key | Desc | Default |
| --- | ---- | ------- |
| CLOCK_WHO | 称呼 | Coder |
| CLOCK_BIRTHDAY | 来都来了 | 1996-04-12 |
| CLOCK_PASS_AWAY | 走就走吧 | 2086-04-12 |
| CLOCK_START_COLOR | 渐变起始颜色 | #F3C8ED |
| CLOCK_END_COLOR | 渐变终止颜色 | #B2F6EF |

新增 .life_clock 文件，以 yaml 格式配置，读取的优先级如下：
1. (高)读取 clock 程序运行路径下的 .life_clock 文件 
2. 读取 $HOME/.life_clock 文件
3. 读取 /etc/.life_clock 文件

### 🤔 人生阿！
> 真正的勇士敢于直面赤裸裸的人生以及血淋淋的事实

![image](https://user-images.githubusercontent.com/19553554/143389082-c135142f-d493-41a8-8774-529b9554dabf.png)

### Contributors

- [NickChenyx](https://github.com/nickChenyx)

### 🔖 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
