rm -rf crypto-config
cryptogen generate --config=crypto-config.yaml
export FABRIC_CFG_PATH=$PWD
rm -rf channel-artifacts
mkdir channel-artifacts
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID byfn-sys-channel
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID wulusaichannel
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID wulusaichannel -asOrg Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID wulusaichannel -asOrg Org2MSP
#docker-compose down --volumes --remove-orphans
#peer channel create -o orderer.wulusai.net:7050 -c wulusaichannel -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/wulusai.net/tlsca/tlsca.wulusai.net-cert.pem