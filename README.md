# 一、gops
    gin练手项目

# 二、初始化配置
```
#初始化数据库表结构
curl -X POST http://127.0.0.1:9000/InitTable

{"code":20000,"message":"初始化表成功"}

#插入默认用户账号信息(admin/admin)
insert into User(`id`, `username`, `nick_name`, `email`, `phone`, `password`, `created_at`) values (1, "admin", "系统管理员", "admin@qq.com", "1234567890", "ktDRP8O2qn2PFLV0yBiHGA==", 0);
```

# 三、postman测试
```
1、#获取token
POST  http://localhost:9000/UserLogin
{
	"username": "admin",
	"password": "admin"
}

2、#设置header
X-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwibmlja25hbWUiOiLns7vnu5_nrqHnkIblkZgiLCJleHAiOjE1NjU2ODIwMjUsImlzcyI6Imp3dC1nbyJ9.8eRsFNf_21FdI7N-EhGZFeCe0HLQyw3zGFKL2tFj7kg

3、#获取用户列表
GET  http://localhost:9000/sys/UserList?page_size=100&current_page=1

4、#新增用户
POST  http://localhost:9000/sys/UserAdd

{
    "username": "cese",
    "password": "123456",
    "nick_name": "测试",
    "email": "1151980610@qq.com",
    "phone": "1832091411"
}

5、#删除用户
POST  http://localhost:9000/sys/UserDelete
{
    "id": 2
}

6、#根据token解析用户信息和权限
GET  http://localhost:9000/sys/AccountInfo


7、#查询用户表，获取所有用户的中文名字和英文名字
GET  http://localhost:9000/sys/UserOptions

测试结果
{
    "code": 20000,
    "data": [
        {
            "nickname": "系统管理员",
            "username": "admin"
        },
        {
            "nickname": "测试1",
            "username": "cese1"
        },
        {
            "nickname": "测试2",
            "username": "cese2"
        }
    ],
    "message": "请求成功"
}
```