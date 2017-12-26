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
	GetNotifiers() []NotifierType
	Fetch() ([]*Item, error)
}
```

定义一个 struct 实现如上4个接口方法

1. 名字
2. 更新间隔
3. 需要怎么通知
4. fetch 的时候去请求并组装好 Item 对象

> 可以看 /commands 里面的一些实例

# TODO
* 清理一些对用户没用的东西
	* init_mongodb