#! /bin/bash
export FABRIC_CFG_PATH=$PWD
if [ $1 == "help" ];then
  echo "fabric development create help:"
  echo " clean  清理配置文件"
  echo " config 生成文件"
  echo " up     打开容器"
  echo " down   关闭容器"
  echo " restart   清理环境重启容器"
fi
function clean() {
  rm -rf crypto-config
  rm -rf channel-artifacts
  echo "clean config  $?"
}
#清理配置文件
if [ $1 == "clean" ];then
  clean
fi

function config() {
  cryptogen generate --config=crypto-config.yaml
  mkdir channel-artifacts
  configtxgen -profile SampleDevModeKafka -outputBlock ./channel-artifacts/genesis.block -channelID byfn-sys-channel
  configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID wulusaichannel
  configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID wulusaichannel -asOrg Org1MSP
  configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID wulusaichannel -asOrg Org2MSP
  echo "create config $?"
}

  #生成文件
if [ $1 == "config" ];then
  config
fi

function up() {
    docker-compose up -d
    echo " docker-compose up $?"
}
#打开容器
if [ $1 == "up" ];then
    up
fi
function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
  else
    docker rm -f $CONTAINER_IDS
  fi
}
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

function down() {
    docker-compose down --volumes --remove-orphans
    #docker rm -f $(docker ps -a | grep "hyperledger/*" | awk "{print \$1}")
    docker volume prune
    export PATH=${PWD}/../bin:${PWD}:$PATH
    docker run -v $PWD:/tmp/wulusai --rm hyperledger/fabric-tools:$IMAGETAG rm -Rf /tmp/wulusai/ledgers-backup
    clearContainers
    removeUnwantedImages
    echo "remove docker-compose $?"
}

#关闭容器
if [ $1 == "down" ];then
    down
fi

function create() {
    peer channel create -o orderer.wulusai.net:7050 -c wulusaichannel -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/tlsca/tlsca.wulusai.net-cert.pem
    echo "create channel block $?"
}

#创建channel.block
if [ $1 == "create" ];then
    create
fi
#
#docker stop $(docker ps -a -q)
#docker rm $(docker ps -a -q)
#docker rmi -f $(docker images |grep "dev-" |awk '{print $3}')
if [ $1 == "restart" ];then
  down
  clean
  config
  up
fi