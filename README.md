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