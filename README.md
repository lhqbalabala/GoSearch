## 该项目的单机架构及单机详细文档
点击 [这里](https://github.com/lgdSearch/lgdSearch) 跳转

## 在线体验
在 2023年之前可以通过 [这里](http://121.196.207.80:8080) 在线体验本项目

## 技术栈

+ Trie树
+ AC自动机
+ 快速排序法
+ 倒排索引
+ 文件分片
+ golang-jieba分词
+ boltdb
+ badgerdb
+ colfer序列化
+ grpc
+ ResNet50
+ diskcache
+ milvus
+ protobuf序列化
+ redis哨兵模式
+ redis分片集群
+ sentinel-go
+ etcd注册中心
+ go-micro
+ 一致性hash
+ apinto网关
+ zookeeper
+ kafka集群

## 该项目与原有单机架构的不同
使用了分布式及微服务架构，实现了水平拆分与垂直拆分

使用etcd作为注册中心，实现服务注册发现

go-micro框架与grpc负责各服务间通信及负载均衡实现

Apinto作为网关，转发相应请求，并安装了流量控制以及额外参数插件，外部请求无法访问Api

若搜索结果大于5W，则用5w大小的堆排代替快排，并把请求放入redis中，结果放入Badgerdb中并设置一小时实现，若搜索请求在redis中有，则直接从本地Badgerdb中取出该搜索结果集合，因为
对于请求设置了一致性hash策略，所以相同请求一般会进入同一Service，从而命中缓存

热搜则将原有堆排转为Redis哨兵集群的Zset，每个Service维护的队列依然不变，进队列时Redis中该值+1，出队列时Redis中该值-1，获取时直接返回Redis的Zset即可

增加了Kafka来集合多个服务器的日志，以及相关搜索的数据增加

使用了Sentinel-go来堆HotSearchRPC限流、熔断、降级，降级策略为若熔断或限流，则直接从本地缓存中取热搜，在每次获得热搜时，本地缓存会实时更新

图片上传中，也将原有本地数据判断图片是否上传改为存于Redis哨兵集群中，保证获取正确的图片是否上传信息

## 本项目当前的服务部署示例图
![服务部署](https://user-images.githubusercontent.com/71314983/184846392-482324c3-c09c-4e93-8cb1-a228ee582143.jpg)

## 执行一次搜索的示例图

![搜索请求](https://user-images.githubusercontent.com/71314983/184942343-bd8950a7-a70c-47cb-9064-2da5f11bd022.jpg)


## 模块功能

UserApi：用户登录、注册、收藏夹等接口实现，并操作MySQL数据库

SearchApi：搜索功能的接口实现

SearchRPC：搜索功能的方法实现

HotSearchApi：热搜的接口实现

HotSearchRPC：热搜的方法实现

RelatedSearchApi：相关搜索的接口实现

RelatedSearchRPC：相关搜索的方法实现

