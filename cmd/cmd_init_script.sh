#!/usr/bin/env bash



# 1. postgres container
mkdir ~/postgresql

docker run --name postgres \
-e POSTGRES_USER=dendrite \
-e POSTGRES_PASSWORD=itsasecret \
-p 5432:5432 \
-v ~/postgresql:/var/lib/postgresql/data \
-d postgres

# 2. pgAdmin4
docker run -p 5050:80 \
    -e "PGADMIN_DEFAULT_EMAIL=user@domain.com" \
    -e "PGADMIN_DEFAULT_PASSWORD=SuperSecret" \
    -d dpage/pgadmin4

# browser: http://localhost:5050  login ->   user@domain.com SuperSecret

# create connect
#hostName: host.docker.internal
#port: 5432
#Maintenance database name: postgres
#userName: dendrite
#password: itsasecret

createdb -U dendrite -O dendrite dendrite_userapi_accounts
createdb -U dendrite -O dendrite dendrite_mediaapi
createdb -U dendrite -O dendrite dendrite_syncapi
createdb -U dendrite -O dendrite dendrite_roomserver
createdb -U dendrite -O dendrite dendrite_keyserver
createdb -U dendrite -O dendrite dendrite_federationapi
createdb -U dendrite -O dendrite dendrite_appservice
createdb -U dendrite -O dendrite dendrite_mscs

# 3. install nats latest
docker run --name nats-server -p 4222:4222 -d nats:latest

# 4. create signing key
cd dendrite
#export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=socks5://127.0.0.1:7890      #use proxy
./build.sh
./bin/generate-keys --private-key matrix_key.pem
cp ./bin/matrix_key.pem ~/

# 5. copy dendrite.yaml and create cert
vim /Users/alan/Desktop/source_read/my_github_code/dendrite/dendrite-sample.polylith.yaml
global:
  private_key: /Users/alan/matrix_key.pem

cp /Users/alan/Desktop/source_read/my_github_code/dendrite/dendrite-sample.polylith.yaml ~/dendrite.yaml

# 6. add run arguments
-config /Users/alan/dendrite.yaml appservice



# 调优
a. 每个服务数据库连接数
b. nats 集群模式



