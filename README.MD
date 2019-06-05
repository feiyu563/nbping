nbping 简介
-----------------

nbping是为解决局域网大批量IP实例或主机探活,采用go协程并发处理,可以自定义并发的协程数量和输出结果.效率远高于现有的批量ping工具.

--------------------------------------

nbping具备如下特性
---------------------
 - 支持结果导出到文件
 - 支持高并发,默认开启20个协程进行处理
 - 支持指定文件批量读取探活IP

--------------------------------------
部署方式
----

nbping可以部署在本地和云平台上，支持windows、linux系统：

 - 在linux系统中部署
```
git clone https://github.com/feiyu563/nbping.git
cd nbping/cmd
./nbping -i ip.txt -o out.csv #如需后台运行请执行nohup ./nbping -i ip.txt -o out.csv &
```
- 在windows系统中运行
```
git clone https://github.com/feiyu563/nbping.git
#打开windows cmd窗口
cd nbping/cmd
nbping.exe -i ip.txt -o out.csv
```


配置说明
----
--------------------------------------

nbping 支持以下参数:

```
  -d int
        if open debug,是否开启debug模式,开启则显示每条探活记录,默认不开启,支持参数0,1,此参数会影响性能,建议大批量时候不开启
  -h    
        for help 显示帮助信息
  -i string
        ip file path (default "ip.txt"),指定IP文件存放的路径(必须参数)
  -n int
        number of thread for ping (default 20),指定同时开启的协程数量,默认20,如果机器性能很高,则可以开启更多数量.
  -o string
        out put file path (default "out.csv"),指定输出的文件路径(必须参数)
  -r int
        if open debug (default 2),失败重试次数.默认为两次重试
```

效果图
----
![tp](https://raw.githubusercontent.com/feiyu563/nbping/master/doc/tp.png)


项目源码
----

 - [nbping][1]


  [1]: https://github.com/feiyu563/nbping

FOR HELP
----
Email: 244217140@qq.com

![me](https://raw.githubusercontent.com/feiyu563/PrometheusAlert/master/doc/me.jpg)