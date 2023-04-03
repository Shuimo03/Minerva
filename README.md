# Minerva

因为K8S event的有效保存时间只有1小时，而event可以有效的分析集群，所以收集event并保存event是一种很重要的事情，本项目目前为练手项目，主要是为了熟悉
client-go的使用，以及收集event并最后输出到类似于es,redis等等地方。

## 项目参考
主要是参考了阿里云的kube-eventer，以及kubesphere的kube_events_exporter:
+ [kube-eventer](https://github.com/AliyunContainerService/kube-eventer)
+ [kube-events](https://github.com/kubesphere/kube-events/)

### timeline
+ 2023/4/3 代码跑通，大致逻辑实现了，目前支持输出到控制台。