# realworld开发日志

## 2025- 5-8号—–分支（master）

1. 查阅相关对应的开发文档内容

2. [Real World](https://github.com/gothinkster/realworld) 是一个相对完整的CRUD实践练习项目. 该项目规定了API 的标准和规范, 开发人员根据自己熟悉的技术栈, 选择实现前端或者后端.


   - 项目Feature 介绍: https://realworld-docs.netlify.app/implementation-creation/introduction/

   - 后端接口文档: https://realworld-docs.netlify.app/specifications/backend/introduction/

3. 制定开发顺序，明确各个功能关系

   ```
   1.用户表
   2.文章表
   3.关注关系中间表 用户id=>用户id
   4.文章喜爱表    用户id=>文章id
   5.评论表       用户id 文章id  评论id
   6.标签表         标签id
   ```

   其中，用户表为优先设计，密码暂时选用明文形式，方便查阅数据库和接口测试  后续采用MD5+随机种子加密

   文章表需包含以下字段，作者id ，文章摘要，文章描述和文章主题，创建时间和更新时间，文章标签

   其中文章标签和标签表无特殊关系，大致流程应该是，前端读取数据库中标签表的数据，渲染到前端按钮，用户点击按钮将标签加入文章标签中，将多个标签结果存放到json格式中，返回到后端，后端直接存放标签数组即可，因此无需绑定tag表，**同时标签表 更新频率低 tag可以在服务启动前，先将标签信息存放到redis中**

   关注关系中间表，根据用户id和用户id进行绑定，将这两个关键字段合并为一个i唯一index，防止重复录入

   评论表，一个用户可以评论多条同样数据的信息，因此不做限定，制作外键绑定

   

   用户表功能

   ```
   1.登录
   2.注册
   3.获取用户数据
   4.更新用户  
   5.删除用户
   ```

   用户表+关注关系表

   ```
   1.关注用户+返回数据
   1.取消关注用户+返回数据
   ```

   文章表+（用户表+关注关系表）(作者信息就是上述返回的结构体)+文章喜爱表

   ```
   1.列出文章
   2.提要文章
   3.获取文章
   4.创建文章
   5.更新文章
   6.设置喜欢
   7.设置不喜欢
   ```

   用户表+文章表

   ```
   删除文章
   ```

   用户表+文章表+评论表+关注关系表（这里的关系是根据用户id，文章id 将信息保存到评论表中，但是返回数据的时候需要返回该评论与用户的关联关系）

   ```
   1.增加评论
   2.获取文章评论
   ```

   

   用户表+文章表+评论表

   ```
   1.删除评论
   ```

   标签表

   ```
   获取所有标签   
   ```

   

   开发配置

   ![](./image/%E5%BC%80%E5%8F%91%E9%85%8D%E7%BD%AE.png)

   2.框架搭建

![](./image/%E6%A1%86%E6%9E%B6%E6%90%AD%E5%BB%BA.png)









## 2025- 5-9号—–分支（master）

1.开发构建模型，分为请求体、响应体和模型生成表

![](./image/%E6%A8%A1%E5%9E%8B%E6%90%AD%E5%BB%BA.png)

2.路由设置 按照功能进行划分句柄，将需要验证和不需要验证进行分组

![](./image/%E5%8A%9F%E8%83%BD%E5%88%92%E5%88%86.png)

![](./image/%E5%88%86%E7%BB%84.png)

3.实现用户登录和注册，以及用户登录和注册

```
1.登录请求体
{
  "user":{
    "email": "jake@jake.jake",
    "password": "jakejake"
  }
}


2.注册请求体
{
  "user":{
    "username": "Jacob",
    "email": "jake@jake.jake",
    "password": "jakejake"
  }
}
```

```
1.响应体
{
  "user": {
    "email": "jake@jake.jake",
    "token": "jwt.token.here",
    "username": "jake",
    "bio": "I work at statefarm",
    "image": null
  }
}
```

单元接口测试：

```
[
  {
    "EMAIL": "user1@example.com",
    "PASSWORD": "password1"
  },
  {
    "EMAIL": "user2@example.com",
    "PASSWORD": "password2"
  },
  {
    "EMAIL": "user3@example.com",
    "PASSWORD": "password3"
  },
  {
    "EMAIL": "user4@example.com",
    "PASSWORD": "password4"
  },
  {
    "EMAIL": "user5@example.com",
    "PASSWORD": "password5"
  }
]


[
  {
    "EMAIL": "user1@example.com",
    "PASSWORD": "password1",
    "USERNAME": "user1"
  },
  {
    "EMAIL": "user2@example.com",
    "PASSWORD": "password2",
    "USERNAME": "user2"
  },
  {
    "EMAIL": "user3@example.com",
    "PASSWORD": "password3",
    "USERNAME": "user3"
  },
  {
    "EMAIL": "user4@example.com",
    "PASSWORD": "password4",
    "USERNAME": "user4"
  },
  {
    "EMAIL": "user5@example.com",
    "PASSWORD": "password5",
    "USERNAME": "user5"
  }
]

```

结果响应

![](./image/%E6%B5%8B%E8%AF%95%E7%BB%93%E6%9E%9C.png)

3.将用户信息保存到token

## 2025- 5-10号—–分支（master）

### 时间太赶，交代开发功能

1.关注功能全部实现

2.文章创建，删除，更新

```
==========实现路由==========
 登录POST   /api/users/login       
 注册POST   /api/users             
获取当前用户 GET    /api/user              
 更新PUT    /api/user              
 创建POST   /api/articles          
 获取单个GET    /api/articles/:slug    
 更新文章PUT    /api/articles/:slug    
 删除文章DELETE /api/articles/:slug    
 
 
 查看关系，增加关系，删除关系
 GET    /api/profiles/:username 
POST   /api/profiles/:username/follow
DELETE /api/profiles/:username/follow
```





## 2025- 5-11号—–分支（master）

### 完成余下功能剩余标签表的生成和文章提要

完成

```
创建				  	 POST   /api/articles          
获取单个				GET    /api/articles/:slug       
更新单个				PUT    /api/articles/:slug       
删除单个				DELETE /api/articles/:slug       
条件获取				GET    /api/articles             
查看评论				GET    /api/articles/:slug/comments 
增加评论				POST   /api/articles/:slug/comments 
删除评论				DELETE /api/articles/:slug/comments/:id
设置喜欢				POST   /api/articles/:slug/favorite 
取消喜欢				DELETE /api/articles/:slug/favorite 
```

