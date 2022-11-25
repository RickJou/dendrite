#!/usr/bin/env bash

# 1. postgres container
```bash
mkdir -p /opt/sola/postgresql
docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=abcd.1234 -p 5432:5432 -v /opt/sola/postgresql:/var/lib/postgresql/data -d postgres

#进入容器,创建必要的8个库
/usr/bin/createdb -U root -O root dendrite_userapi
/usr/bin/createdb -U root -O root dendrite_mediaapi
/usr/bin/createdb -U root -O root dendrite_syncapi
/usr/bin/createdb -U root -O root dendrite_roomserver
/usr/bin/createdb -U root -O root dendrite_keyserver
/usr/bin/createdb -U root -O root dendrite_federationapi
/usr/bin/createdb -U root -O root dendrite_appservice
/usr/bin/createdb -U root -O root dendrite_mscs
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

# 3. install nats latest
mkdir -p /opt/sola/nats/data     #jetstream数据存储目录
mkdir -p /opt/sola/nats/config   #配置文件目录

vim /opt/sola/nats/config/nats-server.conf
```bash
server_name: "sola_nats"
port: 4222
monitor_port: 8222
jetstream: enabled
jetstream {
    store_dir: /data/jetstream
    max_mem: 1G
    max_file: 100G
}
```
# start nats container
```bash
docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 \
-v /opt/sola/nats/config/nats-server.conf:/etc/nats/config/nats-server.conf \
-v /opt/sola/nats/data:/data/jetstream \
nats:2.9.7-linux -c /etc/nats/config/nats-server.conf
```

# 4. create signing key
```bash
cd dendrite
#export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7890      #use proxy
./build.sh
./bin/generate-keys --private-key matrix_key.pem
cp ./matrix_key.pem /opt/sola
```

# 5. copy dendrite.yaml and create cert
```bash
cp /Users/alan/Desktop/source_read/my_github_code/dendrite/dendrite-sample.polylith.yaml /opt/sola/dendrite.yaml
```

# 6. go get 获取项目依赖
go env -w  GOPROXY=https://goproxy.cn,direct
go get

# 7. 启动8个服务(建议: idea 6GB内存)
```sh
dendrite/cmd/dendrite-sola/xxx-xxx-main/main.go              
```

# 8. 启动nginx container 反向代理,APP通过nginx代理

vim /opt/sola/sola.conf   #将三处 192.168.31.194 修改为自己局域网ip,方便移动端请求.
```shell
server {
    #listen 443 ssl; # IPv4
    #listen [::]:443 ssl; # IPv6
    
    listen 8088; # IPv4
    #listen [::]:8088; # IPv6
    server_name 192.168.31.194;

    #ssl_certificate /path/to/fullchain.pem;
    #ssl_certificate_key /path/to/privkey.pem;
    #ssl_dhparam /path/to/ssl-dhparams.pem;

    proxy_set_header Host      $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_read_timeout         600;

    location /.well-known/matrix/server {
        return 200 '{ "m.server": "192.168.31.194:8088" }';
    }

    location /.well-known/matrix/client {
        return 200 '{ "m.homeserver": { "base_url": "http://192.168.31.194:8088" } }';
    }

    # route requests to:
    # /_matrix/client/.*/sync
    # /_matrix/client/.*/user/{userId}/filter
    # /_matrix/client/.*/user/{userId}/filter/{filterID}
    # /_matrix/client/.*/keys/changes
    # /_matrix/client/.*/rooms/{roomId}/messages
    # /_matrix/client/.*/rooms/{roomId}/context/{eventID}
    # /_matrix/client/.*/rooms/{roomId}/event/{eventID}
    # /_matrix/client/.*/rooms/{roomId}/relations/{eventID}
    # /_matrix/client/.*/rooms/{roomId}/relations/{eventID}/{relType}
    # /_matrix/client/.*/rooms/{roomId}/relations/{eventID}/{relType}/{eventType}
    # /_matrix/client/.*/rooms/{roomId}/members
    # /_matrix/client/.*/rooms/{roomId}/joined_members
    # to sync_api
    location ~ /_matrix/client/.*?/(sync|user/.*?/filter/?.*|keys/changes|rooms/.*?/(messages|.*?_?members|context/.*?|relations/.*?|event/.*?))$  {
        proxy_pass http://host.docker.internal:7770;
    }
    #client_api
    location /_matrix/client {
        proxy_pass http://host.docker.internal:8071;
    }
    #federation_api
    location /_matrix/federation {
        proxy_pass http://host.docker.internal:8072;
    }
    #federation_api
    location /_matrix/key {
        proxy_pass http://host.docker.internal:8072;
    }
    #media_api
    location /_matrix/media {
        proxy_pass http://host.docker.internal:8074;
    }
}
```

```shell
docker run --name sola-nginx -p 8088:8088 -p 80:80 -v /opt/sola/sola.conf:/etc/nginx/conf.d/sola.conf -d nginx

#nginx 容器内
curl http://host.docker.internal:8071/_matrix/client/versions
#宿主机
curl http://192.168.31.194:8088/_matrix/client/versions

curl http://192.168.31.194:8088/.well-known/matrix/client

#fluffyChat APP 连接

```





# 调优
a. 每个服务数据库连接数
b. nats 集群模式优化
c. 断开matrix联盟


