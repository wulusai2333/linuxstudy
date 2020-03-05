#!/bin/bash
# Environment variables for PEER0
# peer0.org1
if [ $1 -eq 1 ];then
  export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer0.org1.wulusai.net/tls/server.crt
  export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer0.org1.wulusai.net/tls/server.key
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/users/Admin@org1.wulusai.net/msp
  export CORE_PEER_ADDRESS=peer0.org1.wulusai.net:7051
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer0.org1.wulusai.net/tls/ca.crt
  echo "set peer0.org1 development sucess"
fi
# peer1.org1
if [ $1 -eq 2 ];then
  export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer1.org1.wulusai.net/tls/server.crt
  export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer1.org1.wulusai.net/tls/server.key
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/users/Admin@org1.wulusai.net/msp
  export CORE_PEER_ADDRESS=peer1.org1.wulusai.net:8051
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.wulusai.net/peers/peer1.org1.wulusai.net/tls/ca.crt
  echo "set peer1.org1 development sucess"
fi
# peer0.org2
if [ $1 -eq 3 ];then
  export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/server.crt
  export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/server.key
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/users/Admin@org2.wulusai.net/msp
  export CORE_PEER_ADDRESS=peer0.org2.wulusai.net:9051
  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer0.org2.wulusai.net/tls/ca.crt
  echo "set peer0.org2 development sucess"
fi
# peer1.org2
if [ $1 -eq 4 ];then
  export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer1.org2.wulusai.net/tls/server.crt
  export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer1.org2.wulusai.net/tls/server.key
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/users/Admin@org2.wulusai.net/msp
  export CORE_PEER_ADDRESS=peer1.org2.wulusai.net:10051
  export CORE_PEER_LOCALMSPID="Org2MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.wulusai.net/peers/peer1.org2.wulusai.net/tls/ca.crt
  echo "set peer1.org2 development sucess"
fi

