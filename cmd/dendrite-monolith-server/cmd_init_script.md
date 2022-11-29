#!/usr/bin/env bash

# 1. postgres container
```bash
mkdir -p /opt/sola/postgresql
docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=abcd.1234 -p 5432:5432 -v /opt/sola/postgresql:/var/lib/postgresql/data -d postgres

#进入容器,创建必要的8个库
/usr/bin/createdb -U root -O root dendrite
```
# 2. pgAdmin4
```bash
docker run -p 5050:80 -e "PGADMIN_DEFAULT_EMAIL=root@sola.com" -e "PGADMIN_DEFAULT_PASSWORD=abcd.1234" -d --name pgadmin4 dpage/pgadmin4

# browser: http://localhost:5050  login ->   root@sola.com abcd.1234
# create connect
#hostName: host.docker.internal
#port: 5432
#Maintenance database name: postgres
#userName: dendrite
#password: itsasecret
```

# 3. create signing key
```bash
cd dendrite
#export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7890      #use proxy
./build.sh
./bin/generate-keys --private-key matrix_key.pem
cp ./matrix_key.pem /opt/sola
```

# 4. copy dendrite.yaml and create cert
*注意: 将servername改成自己的ip:nginx端口 -->  192.168.1.120:8088*
//*注意: 开发环境federation禁用TLS   disable_tls_validation: true*
```bash
cp /Users/alan/Desktop/source_read/my_github_code/dendrite/dendrite-sample.monolith.yaml /opt/sola/dendrite.yaml
```

# 5. go get 获取项目依赖
go env -w  GOPROXY=https://goproxy.cn,direct
go get

# 6. 启动服务(建议: idea 6GB内存)
```sh
dendrite/cmd/dendrite-monolith-server/main.go
```

# 7. 部署nginx,做反向代理和配置wellknow端点
```shell
vim /opt/sola/sola.conf   #将三处 192.168.1.120 修改为自己局域网ip,方便移动端请求.
#vim /opt/sola/sola.conf   #将三处 192.168.31.194 修改为自己局域网ip,方便移动端请求.


server {
    listen 8088; # IPv4
    server_name 192.168.1.120;

    proxy_set_header Host      $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_read_timeout         600;

    location /.well-known/matrix/server {
        return 200 '{ "m.server": "192.168.1.120:8088" }';
    }

    location /.well-known/matrix/client {
        return 200 '{ "m.homeserver": { "base_url": "http://192.168.1.120:8088" } }';
    }

    #media_api
    location / {
        proxy_pass http://host.docker.internal:8008;
    }
}
docker run --name sola-nginx -p 8088:8088 -v /opt/sola/sola.conf:/etc/nginx/conf.d/sola.conf -d nginx
```

# 网络测试
```shell
#nginx 容器内
curl http://host.docker.internal:8088/.well-known/matrix/client
curl http://host.docker.internal:8088/_matrix/client/versions
#宿主机
curl http://192.168.1.120:8088/.well-known/matrix/client
curl http://192.168.1.120:8088/_matrix/client/versions
#fluffyChat APP 连接
http://192.168.1.120:8088
```
### 操作
```shell
#服务器地址
http://192.168.1.120:8088
#登录
alan5 alan.1234
alan6 alan.1234
#邀请创建新的聊天室
@alan5:192.168.1.120:8088
@alan6:192.168.1.120:8088
```

# 调优
a. 每个服务数据库连接数
b. nats 集群模式优化
c. 断开matrix联盟


