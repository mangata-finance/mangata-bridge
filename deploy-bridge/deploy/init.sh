#!/bin/bash

set -e

# config_init

# Basepath for mangata-bridge
pushd ../../
BASE_PATH="$PWD"
popd

# ETH truffle path
ETH_TRUFFLE_PATH="$BASE_PATH//ethereum"

# Path to deploy-bridge
export BRIDGE_DEPLOY_PATH="$BASE_PATH/deploy-bridge"


# setup for docker-compose 
yq e ".services.parachain.ports = $SUB_CHAIN_PORTS | .services.parachain.ports[0] style=\"double\"" -i ../docker-compose.yml
yq e ".services.parachain.volumes[0]=\"$SUB_DIR:/var/docker/mangata-custom-node\"" -i ../docker-compose.yml
yq e ".services.parachain.volumes[1]=\"$BRIDGE_DEPLOY_PATH/build/mangataSpec.json:/var/docker/config/mangataSpec.json\"" -i ../docker-compose.yml
yq e ".services.parachain.command=\"$SUB_COMMAND\"" -i ../docker-compose.yml


# Install relevant packages
yarn global add truffle

pushd ../../ethereum
yarn install
popd

pushd ../
yarn install
popd

# Build docker image (using the DockerFile) required to deploy a parachain container locally
docker build -f ./Dockerfile -t mangata/mangata-node .

# Transfer 1 eth from the 0th eth account (seed account) to 1st derived account this is required for tests
# Also dump the private keys and addreses of these two accounts
node -e "require(\"./ethAccountsSetup\").transferEth( '$MGA_Kovan_Bridge_MNEMONIC', '$MGA_Kovan_Bridge_INFURA_ENDPOINT_WS').then(process.exit);"

# Set the 0th eth account as the one used by the relayer
ethPrivateKey0=$(node -e "require(\"./getPrivateKey\").getPrivateKeyFromSeed( '$MGA_Kovan_Bridge_MNEMONIC', '0').then((x)=>{process.stdout.write(x);}).then(process.exit);")
yq e ".services.relayer.environment.ARTEMIS_ETHEREUM_KEY=\"$ethPrivateKey0\"" -i ../docker-compose.yml

# Clean the build folder
pushd ../

sudo rm -rf build

mkdir build
mkdir build/parachain-state
mkdir build/relayer-config
touch build/parachain.env

popd

# Begin deployment starting with deploying the truffle contracts
./deploy-truffle.sh
