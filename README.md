# 如何运行
1. 确认 docker 环境已经装好
2. 参照 `.env.example` 新建 `.env` 文件, 并设置 proxy
	* 如果运行 docker 的环境可以翻墙则跳过这一步
3. 参照 `config.yaml.example` 新建 `config.yaml`
	* 填好 BearyChat incoming 机器人的 url
4. `docker-compose up lili`

# 如何新增一个爬虫


```go
type CommandType interface {
	GetName() string
	GetInterval() time.Duration
	Fetch() ([]*Item, error)
}
```

定义一个 struct 实现如上4个接口方法

1. 名字
2. 更新间隔
3. fetch 的时候去请求并组装好 Item 对象 // 这里主要是推送的内容以及推送到哪里

> 可以看 /commands 里面的一些实例

# Caveats
* 如果 notify 失败的时候暂时不会重试