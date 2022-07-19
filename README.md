> 这是muke上的一个云存储系统项目，计划每天复现一点点。

## Content
- &nbsp;&nbsp;&nbsp;&nbsp;[01 基于命令行的图书的增删查改](https://www.cnblogs.com/cenjw/p/gobeginner-proj-bookstore-cli.html")
-  &nbsp;&nbsp;&nbsp;&nbsp;[02 文件整理](https://www.cnblogs.com/cenjw/p/gobeginner-proj-organize-folder.html)
- &nbsp;&nbsp;&nbsp;&nbsp;[03 Bookstore REST API](https://www.cnblogs.com/cenjw/p/bookstore-rest-api.html)
- &nbsp;&nbsp;&nbsp;&nbsp;[04 Golang仿云盘项目](https://www.cnblogs.com/cenjw/p/go-filestore-disk-system.html)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	- 2.1基础版文件上传  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	- 2.2 文件查询信息接口  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	- 2.3 实现文件下载、修改、删除接口  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;    - 3.1 MySQL主从数据同步(一)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;    - 3.2 持久化  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;     - 4 账号系统和鉴权
- &nbsp;&nbsp;&nbsp;&nbsp;持续更新中...  

# Getting Started

# 准备

- Linux（Ubuntu）/Windows10
- VS Code
- MySQL/Redis/RabbitMQ（同步 to 异步）
- Postman, Chrome
- 云概念（公有云、私有云）

## 要求基础

- Golang基础语法、开发包，有项目开发经验更佳
- 对文件传输和存储场景有兴趣

## 课程安排

- 2-6 构建一个基础版的文件上传服务
- 7-11 架构逐步升级，搭建一个完整优化的分布式服务

## 目标

- 基于 Golang 实现分布式文件上传服务
- 重点结合开源存储（Ceph）及公有云（阿里OSS）支持断点续传及秒传功能
- 微服务化及容器化部署

## ⭐收获

### 工具

- Redis/RabbitMQ
- Docker/Kubernets(k8s)
- 分布式对象存储(Ceph)
- 阿里云OSS对象存储服务

### 干货

- 文件分块断点上传 & 秒传
- 对象从Ceph迁移到阿里云OSS的经验
