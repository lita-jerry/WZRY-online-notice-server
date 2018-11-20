# WZRY-online-notice-server
王者荣耀上线提醒 服务端

## 语言工具
Golang + MySQL

## 实现功能
轮询玩家在线状态, 在线状态一旦改变, 通过`APP push`或`Socket`实现提醒, 并记录在数据库中, 实现统计、分析的功能

## TODO
- [ ] 实现轮询功能, 包含接口请求、响应解析
- [ ] 实现状态变更后的持久, 用户状态变更需要保存在数据库中
- [ ] 实现推送功能, 计划使用[信鸽推送](https://xg.qq.com/), 已支持[服务端Golang接口](https://github.com/xingePush/xinge-api-Golang)
- [ ] 实现`Socket`连接, 更快速推送状态变更消息