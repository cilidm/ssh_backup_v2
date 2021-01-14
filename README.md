### start
> 配置conf.toml,目前支持
>
> 本地->ssh => 1 
>
> 本地->oss store => 2 
>
> ssh->ssh => 4 
>
> ssh->本地 => 6 
>
> oss store->本地 => 7
>
> (oss store目前只测试了七牛云 其他还未测试 如果有bug请提交issue反馈)
>

###配置完毕后执行
`go run main.go -r -s -c`
>
> -r 开启传输任务
>
> -s 开启http服务
>
> -c 清除缓存数据，如要多次执行同路径的任务，请加上此配置项
>

####测试环境 mac linux (win部分适配完毕，如还有问题请提交issue)

###交流群 37703811 