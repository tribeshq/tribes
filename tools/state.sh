#!/bin/bash

# Trap SIGINT (Ctrl+C) to gracefully terminate all child processes
cleanup() {
    echo "Terminating child processes..."
    kill 0 
    exit 0
}
trap cleanup SIGINT

# Helper function to convert string to hex
stringToHex() {
    echo -n "$1" | xxd -p | tr -d '\n' | sed 's/^/0x/'
}

# Ethereum addresses
INPUT_BOX="0xB6b39Fb3dD926A9e3FBc7A129540eEbeA3016a6c"
DAPP_ADDRESS="0x047b88c0fed7e2ce3848cc22dd3406d6de2454dd"
PORTAL_ADDRESS="0x05355c2F9bA566c06199DEb17212c3B78C1A3C31"
ADMIN_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
ADMIN_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
CREATOR_ADDRESS="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
CREATOR_PRIVATE_KEY="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

INVESTOR_ADDRESSES=(
    "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
    "0x90F79bf6EB2c4f870365E785982E1f101E93b906"
    "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
    "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"
    "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
)

INVESTOR_PRIVATE_KEYS=(
    "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
    "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"
    "0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a"
    "0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba"
    "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e"
)

SLEEP_TIME=${SLEEP_TIME:-5}

# Deploy the Token contract and capture the deployed address
deployToken() {
    local tokenName="$1"
    local tokenSymbol="$2"    
    result=$(forge create ./src/Token.sol:Token \
        --private-key $ADMIN_PRIVATE_KEY \
        --rpc-url http://localhost:8080/anvil \
        --root ./contracts \
        --broadcast \
        --constructor-args "$tokenName" "$tokenSymbol" 2>&1)
    deployedAddress=$(echo "$result" | grep "Deployed to:" | awk '{print $3}')
    if [[ -z "$deployedAddress" ]]; then
        echo "Error: Failed to deploy contract for $tokenName ($tokenSymbol)."
        echo "$result"
        exit 1
    fi
    echo "$deployedAddress"
}

# Send input to the INPUT_BOX contract
sendInput() {
    local payload="$1"
    hexPayload=$(stringToHex "$payload")
    cast send $INPUT_BOX "addInput(address,bytes)(bytes32)" $DAPP_ADDRESS $hexPayload --private-key $ADMIN_PRIVATE_KEY --priority-gas-price 10000 --rpc-url http://localhost:8080/anvil $GAS_FLAG
}

# Mint tokens to a specified address
mintTokens() {
    local tokenAddress="$1"
    local recipient="$2"
    local amount="$3"
    cast send $tokenAddress "mint(address,uint256)" $recipient $amount --private-key $ADMIN_PRIVATE_KEY --priority-gas-price 10000 --rpc-url http://localhost:8080/anvil $GAS_FLAG
    echo "Minted $amount tokens to $recipient on $tokenAddress"
}

# Approve ERC20 tokens
approveTokens() {
    local token="$1"
    local spender="$2"
    local amount="$3"
    local privateKey="$4"
    echo "Approving $amount tokens for spender ($spender)..."
    cast send $token \
        "approve(address,uint256)" \
        $spender $amount \
        --private-key $privateKey \
        --priority-gas-price 10000 \
        --rpc-url http://localhost:8080/anvil $GAS_FLAG
}

# Function to deposit ERC20 tokens
depositERC20Tokens() {
    local token="$1"
    local dapp="$2"
    local amount="$3"
    local execLayerData="$4"
    local privateKey="$5"
    echo "Depositing $amount of token ($token) to DApp ($dapp)..."
    cast send $PORTAL_ADDRESS \
        "depositERC20Tokens(address,address,uint256,bytes)" \
        $token $dapp $amount "$(stringToHex $execLayerData)" \
        --private-key $privateKey \
        --priority-gas-price 10000 \
        --rpc-url http://localhost:8080/anvil $GAS_FLAG
}

echo "Deploying contracts..."
STABLECOIN_ADDRESS=$(deployToken "Stablecoin" "STABLECOIN")
sleep $SLEEP_TIME

TOKENIZED_RECEIVABLE_ADDRESS=$(deployToken "Pink" "PINK")
sleep $SLEEP_TIME

echo "Deployed contracts:"
echo "STABLECOIN_ADDRESS=$STABLECOIN_ADDRESS"
echo "TOKENIZED_RECEIVABLE_ADDRESS=$TOKENIZED_RECEIVABLE_ADDRESS"

echo "Minting tokens to investors and creator..."
mintTokens $TOKENIZED_RECEIVABLE_ADDRESS $CREATOR_ADDRESS 10000000
sleep $SLEEP_TIME

mintTokens $STABLECOIN_ADDRESS $CREATOR_ADDRESS 10000000
sleep $SLEEP_TIME

mintTokens $STABLECOIN_ADDRESS ${INVESTOR_ADDRESSES[0]} 10000000
sleep $SLEEP_TIME
mintTokens $STABLECOIN_ADDRESS ${INVESTOR_ADDRESSES[1]} 10000000
sleep $SLEEP_TIME
mintTokens $STABLECOIN_ADDRESS ${INVESTOR_ADDRESSES[2]} 10000000
sleep $SLEEP_TIME
mintTokens $STABLECOIN_ADDRESS ${INVESTOR_ADDRESSES[3]} 10000000
sleep $SLEEP_TIME
mintTokens $STABLECOIN_ADDRESS ${INVESTOR_ADDRESSES[4]} 10000000
sleep $SLEEP_TIME

# Create contracts
echo "Creating contracts..."
sendInput '{"path":"create_contract","payload":{"symbol":"STABLECOIN","address":"'"$STABLECOIN_ADDRESS"'"}}'
sleep $SLEEP_TIME
sendInput '{"path":"create_contract","payload":{"symbol":"TOKENIZED_RECEIVABLE","address":"'"$TOKENIZED_RECEIVABLE_ADDRESS"'"}}'
sleep $SLEEP_TIME

# Create users
echo "Creating users..."
sendInput '{"path":"create_user","payload":{"address":"'"$CREATOR_ADDRESS"'","role":"creator"}}'
sleep $SLEEP_TIME
sendInput '{"path":"create_user","payload":{"address":"'${INVESTOR_ADDRESSES[0]}'","role":"qualified_investor"}}'
sleep $SLEEP_TIME
sendInput '{"path":"create_user","payload":{"address":"'${INVESTOR_ADDRESSES[1]}'","role":"qualified_investor"}}'
sleep $SLEEP_TIME
sendInput '{"path":"create_user","payload":{"address":"'${INVESTOR_ADDRESSES[2]}'","role":"non_qualified_investor"}}'
sleep $SLEEP_TIME
sendInput '{"path":"create_user","payload":{"address":"'${INVESTOR_ADDRESSES[3]}'","role":"non_qualified_investor"}}'
sleep $SLEEP_TIME
sendInput '{"path":"create_user","payload":{"address":"'${INVESTOR_ADDRESSES[4]}'","role":"non_qualified_investor"}}'
sleep $SLEEP_TIME

# Create crowdfunding
echo "Creating crowdfunding..."
current_timestamp=$(date +%s)
closes_at=$((current_timestamp + 80))
maturity_at=$((current_timestamp + 150))
crowdfundingPayload='{"path":"create_crowdfunding","payload":{"max_interest_rate":"10","debt_issued":"100000","fundraising_duration":50,"closes_at":'"$closes_at"',"maturity_at":'"$maturity_at"'}}'
approveTokens $TOKENIZED_RECEIVABLE_ADDRESS $PORTAL_ADDRESS 10000 $CREATOR_PRIVATE_KEY
sleep $SLEEP_TIME # +5s
depositERC20Tokens $TOKENIZED_RECEIVABLE_ADDRESS $DAPP_ADDRESS 10000 "$crowdfundingPayload" $CREATOR_PRIVATE_KEY
sleep $SLEEP_TIME # +5s

# 4. Update crowdfunding to ongoing (sent by admin)
echo "Updating crowdfunding state to 'ongoing'..."
updatePayload='{"path":"update_crowdfunding","payload":{"id":1,"state":"ongoing"}}'
sendInput "$updatePayload" $ADMIN_PRIVATE_KEY

# 5. Create orders from investors (sent by each investor)
echo "Creating orders from investors..."
ORDER_AMOUNTS=(60000 52000 2000 3000 400)
INTEREST_RATES=("9" "8" "4" "6" "4")

orderPayload0='{"path":"create_order","payload":{"crowdfunding_id":1,"interest_rate":"9"}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 60000 ${INVESTOR_PRIVATE_KEYS[0]} &
sleep $SLEEP_TIME # +5s
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 60000 "$orderPayload0" ${INVESTOR_PRIVATE_KEYS[0]} &
sleep $SLEEP_TIME # +5s

orderPayload1='{"path":"create_order","payload":{"crowdfunding_id":1,"interest_rate":"8"}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 52000 ${INVESTOR_PRIVATE_KEYS[1]} &
sleep $SLEEP_TIME # +5s
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 52000 "$orderPayload1" ${INVESTOR_PRIVATE_KEYS[1]} &
sleep $SLEEP_TIME # +5s

orderPayload2='{"path":"create_order","payload":{"crowdfunding_id":1,"interest_rate":"4"}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 2000 ${INVESTOR_PRIVATE_KEYS[2]} &
sleep $SLEEP_TIME # +5s
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 2000 "$orderPayload2" ${INVESTOR_PRIVATE_KEYS[2]} &
sleep $SLEEP_TIME # +5s

orderPayload3='{"path":"create_order","payload":{"crowdfunding_id":1,"interest_rate":"6"}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 3000 ${INVESTOR_PRIVATE_KEYS[3]} &
sleep $SLEEP_TIME # +5s
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 3000 "$orderPayload3" ${INVESTOR_PRIVATE_KEYS[3]} &
sleep $SLEEP_TIME # +5s

orderPayload4='{"path":"create_order","payload":{"crowdfunding_id":1,"interest_rate":"4"}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 400 ${INVESTOR_PRIVATE_KEYS[4]} &
sleep $SLEEP_TIME # +5s
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 400 "$orderPayload4" ${INVESTOR_PRIVATE_KEYS[4]} &
sleep $SLEEP_TIME # +5s

sleep 40
wait

# 7. Close crowdfunding (sent by admin)
echo "Closing crowdfunding..."
closePayload='{"path":"close_crowdfunding","payload":{"creator":"'"$CREATOR_ADDRESS"'"}}'
sendInput "$closePayload" $ADMIN_PRIVATE_KEY

# 9. Settle crowdfunding (sent by creator using stablecoin)
echo "Settling crowdfunding..."
settlePayload='{"path":"settle_crowdfunding","payload":{"crowdfunding_id":1}}'
approveTokens $STABLECOIN_ADDRESS $PORTAL_ADDRESS 108270 $CREATOR_PRIVATE_KEY &
sleep $SLEEP_TIME
depositERC20Tokens $STABLECOIN_ADDRESS $DAPP_ADDRESS 108270 "$settlePayload" $CREATOR_PRIVATE_KEY &
wait