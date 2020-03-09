

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
#疑问?这是在内网多机部署,如果多组织和多节点跨公网和内网该怎么办?
```

#####orderer节点主机的配置

```yaml
#需要注意的是networks的配置
version: '2'

services:
  orderer.wulusai.net: #为了方便看 服务名跟域名相同
    image: hyperledger/fabric-orderer:latest
    container_name: ca.wulusai.net
    environment:
      - FABRIC_LOGGING_SPEC=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0 #orderer节点的监听地址
    volumes:
        - ../crypto-config/...:/var/hyperledger/orderer/msp
        - ../crypto-config/.../tls/:/var/hyperledger/orderer/tls
        - orderer.wulusai.net:/var/hyperledger/production/orderer
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    ports:
      - 7050:7050
    networks:
      default:
      	aliases:
      		- wulusai #此名字是当前配置文件所在的目录的名字
```

##### peer节点主机的配置

```yaml
#需要注意的是节点的networks和extra_hosts的配置
version: '2'

services:

    peer0.org1.wulusai.net:
      container_name: peer0.org1.wulusai.net
      image: hyperledger/fabric-peer:latest
      environment:
        - CORE_PEER_LOCALMSPID=Org1MSP
        - CORE_PEER_ID=peer0.org1.wulusai.net
        - CORE_PEER_ADDRESS=peer0.org1.wulusai.net:7051
        - ...
      volumes:
        - /var/run/:/host/var/run/
        - ...
      working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
      command: peer node start
      networks:
        default:
          aliases:
            - testwork
      ports:
        - 7051:7051
      extra_hosts:  # 声明域名和IP的对应关系 把域名解析为对应IP
        - "orderer.wulusai.net:192.168.100.10"
        #- "peer0.org1.wulusai.net:192.168.100.20"
        
#需要注意的是cli的networks和extra_hosts的配置        
    cli:
      container_name: cli
      image: hyperledger/fabric-tools:latest
      tty: true
      stdin_open: true
      environment:
        - CORE_PEER_ID=cli
        - CORE_PEER_ADDRESS=peer0.org1.wulusai.net:7051
        - CORE_PEER_LOCALMSPID=Org1MSP
        - CORE_PEER_TLS_ENABLED=true
        - ...
      working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
      command: /bin/bash
      volumes:
          - /var/run/:/host/var/run/
         
      depends_on:   # 启动顺序 有用吗?
        - peer0.org1.wulusai.net     
      networks:
          default:
            aliases:
              - wulusai #这里需要注意 order ,peer 都要在同名目录下 这个也要相同,否则找不到
      extra_hosts: #同网络下其他主机地址,这是个 1 orderer 2 peer 的网络 
        - "orderer.wulusai.net:192.168.100.10"
        - "peer0.org1.wulusai.net:192.168.100.20"
        #- "peer0.org2.wulusai.net:192.168.100.30"
```

部署节点

```shell
#准备好事先生成的channel-artifacts crypto-config的文件
#切换到对应主机上 如 pee0.org1的主机 192.168.100.20
#进入主机文件夹 ~/wulusai
#拷问文件 channel-artifacts crypto-config 到目录中
#编写docker-compose.yaml
理解起来很容易 拷贝的文件是工具生成的,里面包含了docker容器启动需要的文件
#容器的启动,cli的操作与之前几乎一样,无非就是容器配置文件被拆分了
容器启动后可执行的操作
#orderer容器: 假设证书文件
	先启动,不做任何操作
#peer0.org1节点: 
	docker-compose up -d 启动后 
	cli:
		执行create channel 
    	peer join
    	peer install操作
    	docker cp cli:/..../peer/channel.block #拷贝通道文件到宿主机中 发给其他主机
    	#还有一种操作方法,将文件放到容器的挂载目录中
#peer0.org2节点:
	
```

#### fabric网络搭建过程

```shell
#1.编写组织信息的配置文件,文件中声明多少个组织,每个组织多少个节点多少用户
crypto-config.yaml
#2.生成创世块文件和通道文件,文件中声明配置组织信息,共识机制,区块生成策略,组织关系的概述
configtx.yaml
```



#### fabric网络组织结构

```
客户端
	链接peer需要用户身份账号信息,就可以连到同组织的peer节点
	客户端发起一笔交易
		会发送到所有参与交易的背书节点上
		参加背书的节点进行模拟交易
		背书节点将处理结果发送给客户端你
		如果提案的结果没问题,客户端将交易提交给orderer节点
		orderer节点将交易打包
		leader节点将打包数据同步到当前组织
		当前组织的提交节点将打包数据写入到区块中
fabric-ca-sever
	可以通过他动态创建用户
	可有可无,因为fabirc是个封闭网络,只是提供了些许灵活性
组织
排序节点
	对交易进行排序
		解决双花问题
	对交易打包
peer节点
	背书节点
		进行模拟交易,将结果返回给客户端
		客户端选择的,指定谁去模拟交易
	提交节点
		将orderer节点打包的数据加入到区块链中
		只要是peer节点,就有提交数据的能力
	主节点
		和orderer排序节点直接通信的节点
			从orderer节点处获取到打包数据
			将数据同步到当前组织的各个节点中
		只能有一个
			可以自己指定
			也可通过fabric框架自主选举
	锚节点
		代表当前组织和其他组织通信的节点
		只能有一个
```



SCP远程拷贝

```shell
scp 要拷贝的文件路径 远程主机用户名@远程主机ip:远程主机目录
scp -r 要拷贝的目录 远程主机用户名@远程主机ip:远程主机目录
scp -r /root/wulusai root@192.168.1.2:/root

#安装同一个链码文件时有可能出现指纹不匹配的问题,可以用如下方式安装
#安装链码 在一个节点上安装完链码,打包然后远程拷贝到其他节点主机上
#将链码打包
peer chaincode package -n bhxycc -v 1.0 -p github.com/chaincode bhxycc.1.0.out
#将链码从容器中拷贝到主机上
docker cp cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/bhxycc.1.0.out ./
#scp发送到远程主机
scp ./bhxycc.1.0.out root@192.168.1.2:/root/wulusai/channel-artifacts

#进入另一个主机的cli容器中
docker exec -it cli bash
#安装链码
peer chaincode install ./channel-artifacts/bhxycc.1.0.out
```

#### kafka多级多节点配置

```shell
#为保证集群的可用性,3台主机
#zookeeper主机1:192.168.100.101
#zookeeper主机2:192.168.100.102
#zookeeper主机3:192.168.100.103
#kafka集群至少4个主机才行
#kafka主机1:192.168.100.201 
#kafka主机2:192.168.100.202
#kafka主机3:192.168.100.203
#kafka主机4:192.168.100.204
# 同样最低3台
#orderer主机3:192.168.100.20
#orderer主机3:192.168.100.21
#orderer主机3:192.168.100.22
# 两个组织一个组织一个peer节点
#peer主机3:192.168.100.30
#peer主机3:192.168.100.40
```

zookeeper配置

```yaml
version: '2'

services:

  zookeeper1:
    container_name: zookeeper1
    hostname: zookeeper1
    image: hyperledger/fabric-zookeeper:latest
    restart: always
    environment:
      # ID在集合中必须是唯一的并且应该有一个值，在1和255之间。
      - ZOO_MY_ID=1
      # server.x=[hostname]:nnnnn[:nnnnn]
      - ZOO_SERVERS=server.1=zookeeper1:2888:3888 server.2=zookeeper2:2888:3888 server.3=zookeeper3:2888:3888
    ports:
      - 2181:2181
      - 2888:2888
      - 3888:3888
    extra_hosts:
      - zookeeper1:192.168.100.101
      - zookeeper2:192.168.100.102
      - zookeeper3:192.168.100.103
      - kafka1:192.168.100.201
      - kafka2:192.168.100.202
      - kafka3:192.168.100.203
      - kafka4:192.168.100.204
```

kafka配置

```yaml

version: '2'

services:

  kafka1:
    container_name: kafka1
    hostname: kafka1
    image: hyperledger/fabric-kafka
    restart: always
    environment:
      # broker.id
      - KAFKA_BROKER_ID=1 
      - KAFKA_MIN_INSYNC_REPLICAS=2 #最小备份数
      - KAFKA_DEFAULT_REPLICATION_FACTOR=3 #默认备份数
      - KAFKA_ZOOKEEPER_CONNECT=zookeeper1:2181,zookeeper2:2181,zookeeper3:2181
      # 100 * 1024 * 1024 B
      - KAFKA_MESSAGE_MAX_BYTES=104857600  #最大信息个头 根据orderer节点打包区块大小设置,orderer默认99M 这里信息包括消息头 所以给100M
      - KAFKA_REPLICA_FETCH_MAX_BYTES=104857600 #配置同上
      - KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false
      - KAFKA_LOG_RETENTION_MS=-1 #记录日志的时间间隔 -1表示不记录
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M #堆内存 默认1G
    ports:
      - 9092:9092
    extra_hosts:
      - zookeeper1:192.168.100.101
      - zookeeper2:192.168.100.102
      - zookeeper3:192.168.100.103
      - kafka1:192.168.100.201
      - kafka2:192.168.100.202
      - kafka3:192.168.100.203
      - kafka4:192.168.100.204
```

orderer配置

```yaml
version: '2'

services:

  orderer0.wulusai.net:
    container_name: orderer0.wulusai.net
    image: hyperledger/fabric-orderer:latest
    environment:
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=kafka_default #注意这里的网络设置
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=7050
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      # enabled TLS
      - ORDERER_GENERAL_TLS_ENABLED=false
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]

      - ORDERER_KAFKA_RETRY_LONGINTERVAL=10s
      - ORDERER_KAFKA_RETRY_LONGTOTAL=100s
      - ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s
      - ORDERER_KAFKA_RETRY_SHORTTOTAL=30s
      - ORDERER_KAFKA_VERBOSE=true
      - ORDERER_KAFKA_BROKERS=[192.168.100.201:9092,192.168.100.202:9092,192.168.100.203:9092,192.168.100.204:9092] #kafka主机的地址
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ./crypto-config/ordererOrganizations/test.com/orderers/orderer0.wulusai.net/msp:/var/hyperledger/orderer/msp
      - ./crypto-config/ordererOrganizations/test.com/orderers/orderer0.wulusai.net/tls/:/var/hyperledger/orderer/tls
    networks:
      default:
        aliases:
         - kafka
    ports:
      - 7050:7050
    extra_hosts:
      - kafka1:192.168.100.201
      - kafka2:192.168.100.202
      - kafka3:192.168.100.203
      - kafka4:192.168.100.204
```

