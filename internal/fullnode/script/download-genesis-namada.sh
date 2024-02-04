CHAIN_ID=$0
NAMADA_NETWORK_CONFIGS_SERVER="https://github.com/anoma/namada-shielded-expedition/releases/download/$CHAIN_ID"

if [ ! -d $CHAIN_HOME/$CHAIN_ID ]; then
    echo "Directory $CHAIN_ID does not exist. Downloading..."
    namada --base-dir $CHAIN_HOME/namada client utils join-network --chain-id "$CHAIN_ID"
    echo "$CHAIN_ID downloaded successfully."
else
    echo "Directory $CHAIN_ID already exists."
fi
