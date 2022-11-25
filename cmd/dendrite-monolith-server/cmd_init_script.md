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
```bash
cp /Users/alan/Desktop/source_read/my_github_code/dendrite/dendrite-sample.monolith.yaml /opt/sola/dendrite.yaml
```

# 5. go get 获取项目依赖
go env -w  GOPROXY=https://goproxy.cn,direct
go get

# 6. 启动8个服务(建议: idea 6GB内存)
```sh
dendrite/cmd/dendrite-monolith-server/main.go              
```

```shell

curl http://host.docker.internal:8071/_matrix/client/versions
#宿主机
curl http://localhost:8008/_matrix/client/versions

curl http://192.168.31.194:8088/.well-known/matrix/client

#fluffyChat APP 连接

```

# 调优
a. 每个服务数据库连接数
b. nats 集群模式优化
c. 断开matrix联盟


