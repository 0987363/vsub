### ss、ssr、v2ray节点管理，支持订阅分享

+ **支持多个订阅**
+ **每个订阅可添加不同节点**
+ **注意：订阅接口无鉴权**
+ **支持导入分析节点**
+ **支持导入订阅节点**
+ **支持通过url自动导入订阅订阅**


#### 暂无web节目，只有api接口

+ **可参考postman**
` https://www.getpostman.com/collections/7c74c65a548ed57b0873 `


#### 数据库

+ **可选：手动创建索引，个人用户无所谓**
```
use vsub
db.user.createIndex({username: 1})
db.share.createIndex({user_id: 1})
db.node.createIndex({user_id: 1})
```

