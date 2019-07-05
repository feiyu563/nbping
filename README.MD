nbping 简介
-----------------

nbping是为解决局域网大批量IP实例或主机探活,采用go协程并发处理,可以自定义并发的协程数量和输出结果.效率远高于现有的批量ping工具.

--------------------------------------

nbping具备如下特性
---------------------
 - 支持结果导出到文件
 - 支持高并发,默认开启20个协程进行处理
 - 支持指定文件批量读取探活IP
 - 支持自动过滤无效IP,空行等

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
        是否打开debug模式,即开启显示每条IP结果记录
  -h    显示此帮助页
        build by zhangjikun@haima.me
        version 1.1
  -i string
        IP文件所在路径,,默认读取当前目录ip.txt (default "ip.txt")
  -n int
        开启的并发协程数量,默认20 (default 20)
  -o string
        结果输出文件所在路径,默认放在当前目录out.csv (default "out.csv")
  -r int
        IP检测失败重试次数,默认2次 (default 2)
```

效果图
----
![tp](doc/tp.png)


项目源码
----

 - [nbping][1]


  [1]: https://github.com/feiyu563/nbping

FOR HELP
----
Email: 244217140@qq.com

![me](doc/wx.png)
