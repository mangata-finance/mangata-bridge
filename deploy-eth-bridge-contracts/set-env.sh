# Seed/security phrase for depolying account
export MGA_Kovan_Bridge_MNEMONIC=""

# INFURA Project Id for Kovan
export KOVAN_INFURA_PROJECT_ID="44a296fce7bf4ad79e21b1c2002d34d1"

export INFURA_PROJECT_ID=$KOVAN_INFURA_PROJECT_ID
export MNEMONIC=$MGA_Kovan_Bridge_MNEMONIC

# Basepath for mangata-bridge
pushd ../
BASE_PATH="$PWD"
popd

# Path to deploy-bridge
export BRIDGE_CONTRACTS_DEPLOY_PATH="$BASE_PATH/deploy-eth-bridge-contracts"