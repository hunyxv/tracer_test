http 链路追踪 demo（OpenTracing + Jaeger）.


Jaeger 是 Uber 开发的一套分布式追踪系统，更多信息请查看[官方文档](https://www.jaegertracing.io/docs/1.20/getting-started/)。

首先启动一个简单 UI 服务：
```shell
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.20
```
效果展示：
![demo](https://i.loli.net/2020/10/29/TxPvwMC34UH1cEm.png)