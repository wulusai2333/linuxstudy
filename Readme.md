

### fabric环境搭建(基于ubuntu)

#### 准备工作

```
官方帮助文档
安装curl
安装docker
安装docker-compose
安装go
安装node.js
安装python 2.7版本以上
git
git clone -b release-1.4 https://github.com/hyperledger/fabric.git
GO111MODULE=on go mod init bhxy
GO111MODULE=on go build
```



#### 1.安装和更新基础软件

```shell
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl git software-properties-common lrzsz -y
#添加阿里的docker镜像仓库
apt-get install docker-ce -y
#用非root用户操作是可能权限不足 当前用户加入docker组中
sudo gpasswd -a ${USER} docker
systemctl restart docker
newgrp - docker #当前用户切换docker群组
sudo docker version
#装docker-compose
sudo apt-get install python-pip -y
sudo pip install docker-compose
sudo docker-compose version
#安装go
wget go安装包
tar zxvf go安装包 -C /use/local
mkdir $HOME/go
vim ~/.bashrc
	export GOROOT=/user/local/go
	export GOPATH=$HOME/go
	export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
#使环境变量生效
source ~/.bashrc
. ~/.bashrc #对当前用户生效
go version
安装mode.js
wget 安装包
#解压 /opt
sudo tar -xvf 安装包 -C /opt
sudo vim /etc/profile
	export NODEJS_HOME=/opt/node
	export PATH=$PATH:$NODEJS_HOME/bin
. /etc/profile #对操作系统生效
node -v

```

#### 2.安装fabric-sample

```shell
#需要注意fabric项目安装位置在GOPATH的src目录下github.com/hyperledger/
#可以使用git clone拉取项目
#或者参考fabric官方文档使用curl 下载 bootstrap.sh
#下载过程中 hyperledger-fabric-ca-amd64-1.4.4.tar.gz hyperledger-fabric-ca-amd64-1.4.4.tar.gz 这两个包下载极慢,可以直接将此包放到fabric-sample中减少下载时间
#拉取镜像
fabric-peer #peer模块镜像
fabric-orderer #order节点
fabric-ccenv #chaincode运行环境库
fabric-tools #工具镜像包含cryptogen,configtxgen
fabric-ca #ca模块
fabric-couchdb #couchdb数据库
fabric-kafka 
fabric-zookeeper
#镜像位置
/var/lib/docker/image
#fabric-sample 的bin目录加载到环境变量中,参考node.js
./bysh.sh generate #生成证书文件
./bysh.sh up #启动
./bysh.sh down#关闭
```

3.一些概念

tls -> ssl  https 都是加密的,节点间需要证书

#### 逻辑架构

```shell
#程序员需要做的
身份管理inentity	->成员服务membership	->注册登录enroliment 属性证书attributes		

账本管理ledger		->|
					共识服务consensus -> 分布式账本 排序服务	P2P协议 背书验证
交易管理transcations->| 		

智能合约smart contract	->链码服务chaincode ->安全容器环境 安全镜像仓库
#程序员需要做的是最前面的一层 智能合约是最简单的
```

##### 成员管理

```
会员注册
	注册成功的一个账户得到的不是用户名密码 是一个证书
	使用证书做身份认证
身份保护

交易审计

内容保密
	可以多条区块链,通过通道区分

```

##### 账本管理

```
区块链
	保存所有交易了记录
世界状态
	数据最新状态
	数据存储在当前节点的数据库中 默认levelDB

```

##### 交易管理

```
部署交易
	部署的是链码,就是给节点安装链码
调用交易
	invoke
```

#### 基础概念

##### 组织

``` 
->社会实体
组织中:
有用户
进行数据处理的节点
put 写入数据到区块链中
get 数据查询
```

##### 节点

```
client 
	进行交易管理(cli node sdk,java sdk)
	cli ->通过linux命令行进行操作,使用的是shell命令对数据进行提交和查询
	node.js ->api实现客户端 提供服务,浏览器查询
	java-> 同上
	go-> 同上
peer
	存储和同步账本数据
	用户通过客户端工具对数据进行put操作,数据写入到一个节点中
	数据同步是fabric框架实现的
orderer
	排序和分发交易
	交易数据线打包再写入到区块中
```

##### 通道

``` 
->QQ群 只有在同一个群中才能看到一个群的消息
consensus server:orderer节点
peer节点加入一个通道就要创一个区块链
```

交易流程

```
1.Application/SDK 充当客户端
2.客户端发起一个提案,给peer节点
3.peer节点预演,得到一个结果
4.peer节点将交易结果发送给客户端
	如果模拟交易失败,流程终止
	成功继续
5.客户端将交易提交给排序节点
6.排序节点对交易打包
7.orderer节点将打包数据发送给peer,peer节点将数据写入搭配区块中
	打包数据的发送,不是实时的
	有设定条件,在配置文件中
背书策略:
	完成一笔交易的过程就是背书
```

思考:

站在普通人的角度看Application/SDK 实际是后台服务器

交易通过 手机APP/浏览器(终端) --发送->  Application/SDK(后台服务器)处理请求 -->pee背书节点-->后台应用服务器(确认验证成功)  --> orderer节点排序打包区块 --> peer主节点存储分发区块

#### fabric核心模块

```
peer 主节点模块,负责存储区块链数据,运行维护链码
orderer 交易打包,排序模块
crytogen 组织和证书生成模块
configtxgen 区块和交易生成模块
configtxlator 区块和交易解析模块 ->解析成json格式
```

#### cryptogen命令生成证书文件

```shell
crptogen showtemplate > crypto-config.yaml #生成模板配置文件
#修改配置文件内容
cryptogen generate --config=crypto-config.yaml #根据指定配置文件生成证书
```

msp 是什么?

账号

​	谁有msp

​		每个节点都有一个msp账号

​		每个用户都有msp账号

#### 创世块文件和通道文件生成

```shell
#需要从fabric-sample/first-network下复制configtx.yaml并修改
#已知Capabilities规则修改会导致创建通道失败
configtxgen --help
	-outputBlock string #输出创世区块文件的路径和名字
	-channelID string #指定channel名字,没有用默认
	-outputCreateChannelTx string #输出通道文件的路径和名字
	-profile string #指定配置文件中的节点
	-outputAnchorPeersUpdate string #更新channel配置信息
	-asOrg string #指定所属的组织名称
#执行这个命令需要configtx.yaml 可以复制fabric-sample/first-network下的
#引用配置文件的参数生成创世区块 后面是生成文件的路径
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID wulusaichannel
#创建目录放区块文件
mkdir channel-artifacts
#生成通道文件
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID wulusaichannel
#锚节点更新文件 这个操作可选 这个操作主要用来想更换锚节点时使用
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID wulusaichannel -asOrg Org1MSP
#---asOrg:锚节点的组织名
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID wulusaichannel -asOrg Org2MSP
```



peer中四种节点的角色和作用

```
anchor 锚节点 组织与其他组织通信的节点
leader 领导节点 组织与orderer通信的节点
commit 提交节点 把数据写入到区块链中
背书节点 模拟交易
```

修改 各节点容器的docker-compose 配置文件

docker-compose

```shell
#启动 守护进程
docker-compose -f docker-compose-cli.yaml up -d 
#所有容器 up状态为up 且有端口映射才算成功
set IMAGE_TAG=latest  
set COMPOSE_PROJECT_NAME=wulusai #设置为空 则 容器网络等于 _byfn
#创建的网络则是当前目录名 + _byfn
docker-compose up -d 
docker-compose -f docker-compose-cli.yaml ps #查看启动状态
#如果容器状态正常就可以进入cli操作了
#删除容器
docker rm `docker ps -aq` -f
```

#### cli容器操作节点

```shell
步骤总结:
	先用orderer的证书创建通道
	设置环境变量peer0.org1 加入通道
	设置环境变量peer0.org2 加入通道
	对应环境下更新锚节点(如果锚节点不变这步可以不要)
	对应环境install链码
	初始化链码 使用orderer证书
	invoke调用 orderer证书使用两个锚节点证书,指定通道
	query查询
示例:
	#创建通道
peer channel create -o orderer地址:7050 -c 通道名 -f 通道文件 --tls true --cafile orderer节点pem证书文件绝对路径
	#加入通道
peer channel join -b wulusaichannel.block	
	#更新锚节点
peer channel update -o orderer节点地址:端口 -c 通道名 -f 锚节点更新文件 --tls true --cafile orderer节点pem格式证书文件
	#安装链码
peer chaincode install -n 链码名字 -v 链码版本 -l 链码语言 -p 链码位置
	#初始化
peer chaincode instantiate -o orderer节点地址:端口 -tls true -cafile orderer节点pem格式证书文件 -C 通道名称 -n 链码名称 -l 链码语言 -v 链码版本 -c 链码init函数调用 -P 背书策略
	#invoke调用
peer chaincode invoke  -n 链码名字 -c '{"Args":["userRegister", "2", "user1"]}' -o orderer节点地址:端口 --tls true --cafile orderer节点pem格式证书文件 -C wulusaichannel --peerAddresses org1背书节点:端口 --tlsRootCertFiles org1根ca.crt --peerAddresses org2背书节点:端口 --tlsRootCertFiles org2根ca.crt
	#query查询
peer chaincode query -C wulusaichannel -n bhxycc -c '{"Args":["queryUser", "2"]}'
```

##### 环境配置

```shell
环境配置都放在scripts/changepath.sh中了,可以指定不同节点的环境变量执行peer channel join
查看当前节点是否加入通道: peer channel list
```



##### 创建通道

```shell
#创建通道
peer channel create -o orderer地址:7050 -c 通道名 -f 通道文件 --tls true --cafile orderer节点pem证书文件绝对路径
#crypto-config/ordererOrganizations/wulusai.net/tlsca/tlsca.wulusai.net-cert.pem这是宿主机的文件,在cli中应该找客户端的文件绝对路径
peer channel create -o orderer.wulusai.net:7050 -c wulusaichannel -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/tlsca/tlsca.wulusai.net-cert.pem
#此处遇到bug 修改configtx.yaml的Capabilities为fabric-sample的默认值解决,以下为错误信息
#Error: got unexpected status: BAD_REQUEST -- error validating channel creation transaction for new channel 'wulusaichannel', could not succesfully apply update to template configuration: error authorizing update: error validating DeltaSet: policy for [Value]  /Channel/Capabilities not satisfied: implicit policy evaluation failed - 0 sub-policies were satisfied, but this policy requires 2 of the 'Admins' sub-policies to be satisfied
#生成了 通道名.block 文件
```

##### 当前节点加入通道

```shell
#当前节点加入通道 这里用org1的管理员添加一次 再用org2的管理员添加一次
peer channel join -b wulusaichannel.block
```

##### 更新锚节点

```shell
#更新锚节点 在configtx.yaml中已经指定了默认锚节点 如果不需要更换锚节点这一步可以不做
peer channel update -o orderer节点地址:端口 -c 通道名 -f 锚节点更新文件 --tls true --cafile orderer节点pem格式证书文件
#更新锚节点 org1
peer channel update -o orderer.wulusai.net:7050 -c wulusaichannel -f ./channel-artifacts/Org1MSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp/tlscacerts/tlsca.wulusai.net-cert.pem
#更新锚节点 org2
peer channel update -o orderer.wulusai.net:7050 -c wulusaichannel -f ./channel-artifacts/Org2MSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp/tlscacerts/tlsca.wulusai.net-cert.pem
```

##### 安装链码

```shell
#想要在哪个节点上安装链码就需要在哪个节点配置下install
# -l 默认为go
peer chaincode install -n 链码名字 -v 链码版本 -l 链码语言 -p 链码位置
#安装链码 -p 必须是链码安装的目录而不是链码,起始为$GOPATH下路径 /opt/gopath
peer chaincode install -n bhxycc -v 1.0 -p github.com/chaincode/bhxy
```

##### 链码初始化

```shell
#init初始化 只需要任意节点初始化一次,数据会自动同步
#链码的初始化 
peer chaincode instantiate -o orderer节点地址:端口 -tls true -cafile orderer节点pem格式证书文件 -C 通道名称 -n 链码名称 -l 链码语言 -v 链码版本 -c 链码init函数调用 -P 背书策略
#初始化
peer chaincode instantiate -o orderer.wulusai.net:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp/tlscacerts/tlsca.wulusai.net-cert.pem -C wulusaichannel -n bhxycc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
#Error: error endorsing chaincode: rpc error: code = Unknown desc = access denied: channel [wulusaichannel] creator org [Org2MSP]
# 背书策略改为member了
peer chaincode instantiate -o orderer.wulusai.net:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp/tlscacerts/tlsca.wulusai.net-cert.pem -C wulusaichannel -n bhxycc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.member','Org2MSP.member')"
#可能是未加入通道 peer channel list查看情况
```

##### invoke调用

```shell
#调用需要向orderer节点发送请求,然后由背书规则背书,结果发送给orderer打包
peer chaincode invoke  -n 链码名字 -c '{"Args":["userRegister", "2", "user1"]}' -o orderer节点地址:端口 --tls true --cafile orderer节点pem格式证书文件 -C wulusaichannel --peerAddresses org1背书节点:端口 --tlsRootCertFiles org1根ca.crt --peerAddresses org2背书节点:端口 --tlsRootCertFiles org2根ca.crt
#invoke调用  需要根据制定背书策略选择背书节点证书 orderer证书
peer chaincode invoke  -n bhxycc -c '{"Args":["userRegister", "2", "user1"]}' -o orderer.wulusai.net:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp/tlscacerts/tlsca.wulusai.net-cert.pem -C wulusaichannel --peerAddresses peer0.org1.wulusai.net:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer0.org1.wulusai.net/tls/ca.crt --peerAddresses peer0.org2.wulusai.net:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/ca.crt

#Error: endorsement failure during invoke. response: status:500 message:"cannot retrieve package for chaincode bhxycc/1.0, error open /var/hyperledger/production/chaincodes/bhxycc.1.0: no such file or directory" 
#在 install命令都执行后出现 多节点部署问题,解决办法:先打包再按照

#按顺序执行 创建通道,加入通道 更新锚节点 安装链码 初始化 query成功 invoke调用失败
#可能原因1,没有在所有对等节点上安装链码
#可能原因2,当前环境变量设置问题,没有在通道内链码调用invoke
#可以使用 peer channel list 检查当前peer的通道
#可以使用 echo $CORE_PEER_ADDRESS 检查当前peer的服务地址
#Error: could not assemble transaction: ProposalResponsePayloads do not match - proposal response: version:1 response:<status:200 > payload:"...

```

##### 查询

```shell
#查询 可以查询 查询不需要经过orderer,只需要向通道内节点请求就行
peer chaincode query -C wulusaichannel -n bhxycc -c '{"Args":["queryUser", "2"]}'
#升级链码 这个代码无用
peer chaincode Upgrade -o orderer.wulusai.net:7050 -C wulusaichannel -n bhxycc -v 1.1 -c '{"Args":["init"]}'
```

https://godoc.org/github.com/hyperledger/fabric/core/chaincode/shim

配置环境的文件的意义

```shell
# core peer msp config path 当前peer节点的admin的msp证书
#peerOrg下的组织org2 的 users 下的 admin用户的msp
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/users/Admin@org2.wulusai.net/msp
#core peer address 
export CORE_PEER_ADDRESS=peer0.org2.wulusai.net:9051
#core peer local msp id 操作的peer节点名字
export CORE_PEER_LOCALMSPID="Org2MSP"
#core peer tls cert file peer节点的证书
#peerOrg下的组织org2 的 peers 的peer0 的tls下的 server.crt
export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/server.crt
#core peer tls key file peer节点的秘钥文件
#peerOrg下的组织org2 的 peers 的peer0 的tls下的 server.key
export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/server.key
#core peer tls rootcert file peer节点的根文件 根证书
#peerOrg下的组织org2 的 peers 的peer0 的tls下的ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/ca.crt
 #orderer节点的ca证书
 #在ordererOrg目录下的 域名/orderers 的 对应orderer节点 
 #的msp的tls cacerts tlsca.wulusai.net-cert.pem
 --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp/tlscacerts/tlsca.wulusai.net-cert.pem
```

#### fabric账号

> 根据PKI规范生成的一组证书和秘钥文件
>
> 作用:
>
> ​	保证记录在区块链中的数据具有不可逆,不可篡改
>
> ​	fabric中每条交易都会加上发起者的标签(签名证书) 同时用发起人的私钥加密
>
> ​	如果交易需要其他组织的节点提供背书功能,name背书节点也会在交易中加入自己的签名
>
> 创建channel

如何寻找对应的账号目录

```shell
#Orderer 启动路径
crypto-config/ordererOrganizations/wulusai.net/orderers/orderer.wulusai.net/msp
#Peer 启动的账号路径
crypto-config/peerOrganizations/org1.wulusai.net/peers/peer0.org1.wulusai.net/msp
#创建 channel的账号路径 crypto是cli容器内目录
crypto/peerOrganizations/org2.wulusai.net/users/Admin@org2.wulusai.net/msp
```



```shell

```

crypto的tree

```

```

#### fabric-ca



```shell
node.js 的api 编程去组织上生成账号
官方建议,一个组织对应一个ca
#fabric-ca的配置 模板文件fabric-sample/base-network/docker-compose.yaml
docker rm `docker ps -aq` -f #删除容器
#如果没有按照npm
yum install npm
#新建一个目录初始化 遇到选项全回车就行最终生成一个 package.json
npm init
npm install -g node-gyp
npm install -g node-pre-gyp 
npm install -g cnpm --registry=https://registry.npm.taobao.org
npm install --save grpc --unsafe-perm #安装成功
#使用nodejs 依赖包
npm install --save fabric-ca-client --unsafe-perm #安装成功
#这两个安装包安装失败
yum -y update gcc
yum -y install gcc+ gcc-c++
npm install --save fabric-client --unsafe-perm
#如果没能解决可以尝试更新nodejs如下
#更新yum源
mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
yum makecache
yum -y update
#如果没有更新则直接下载最新版 添加到/etc/profile
wget https://nodejs.org/dist/v12.16.1/node-v12.16.1-linux-x64.tar.xz
xz -d node-v12.16.1-linux-x64.tar.xz
tar -xf node-v12.16.1-linux-x64.tar
#创建软连接
ln -s ~/node-v12.16.1-linux-x64/bin/node /usr/bin/node
ln -s ~/node-v12.16.1-linux-x64/bin/npm /usr/bin/npm
#删除软连接 
rm -rf /usr/bin/node
#install报错了试试下面的
npm cache verify
npm cache clean
npm cache clean --force
npm i -g npm
grep -ir "sha1-xxxxxxxxxxxxxxxx" ~/.npm
#执行命令
node enrollAdmin.js #创建管理员用户 > hfc-key-store


#hf-key-store
hfc-key-store/
├── 709630d415d64255d1c9cac3483bf3bd044a3ac8c7c942eeae34795e0f4e0f1d-priv
├── 709630d415d64255d1c9cac3483bf3bd044a3ac8c7c942eeae34795e0f4e0f1d-pub
└── admin
#执行创建普通用户
node enrollUser.js #创建普通用户
hfc-key-store/
├── 1f98230937ffcb5057794d89569e1697fd49a23fd182577ca66d37a3cd8149fc-priv
├── 1f98230937ffcb5057794d89569e1697fd49a23fd182577ca66d37a3cd8149fc-pub
├── 709630d415d64255d1c9cac3483bf3bd044a3ac8c7c942eeae34795e0f4e0f1d-priv
├── 709630d415d64255d1c9cac3483bf3bd044a3ac8c7c942eeae34795e0f4e0f1d-pub
├── admin
└── user3
node query.js
```

#### solo多机多节点部署

```shell
n台主机需要创建一个名字相同的工作目录
#主机1:192.168.100.10
mdir ~/wulusai
#主机2:192.168.100.20
mdir ~/wulusai
#主机3:192.168.100.30
mdir ~/wulusai
```
