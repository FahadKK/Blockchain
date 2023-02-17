#!/usr/bin/env bash

generateArtifacts() {
  printHeadline "Generating basic configs" "U1F913"

  printItalics "Generating crypto material for Orderer" "U1F512"
  certsGenerate "$FABLO_NETWORK_ROOT/fabric-config" "crypto-config-orderer.yaml" "peerOrganizations/orderer.example.com" "$FABLO_NETWORK_ROOT/fabric-config/crypto-config/"

  printItalics "Generating crypto material for CountryA" "U1F512"
  certsGenerate "$FABLO_NETWORK_ROOT/fabric-config" "crypto-config-countrya.yaml" "peerOrganizations/countrya.example.com" "$FABLO_NETWORK_ROOT/fabric-config/crypto-config/"

  printItalics "Generating crypto material for CountryB" "U1F512"
  certsGenerate "$FABLO_NETWORK_ROOT/fabric-config" "crypto-config-countryb.yaml" "peerOrganizations/countryb.example.com" "$FABLO_NETWORK_ROOT/fabric-config/crypto-config/"

  printItalics "Generating genesis block for group group1" "U1F3E0"
  genesisBlockCreate "$FABLO_NETWORK_ROOT/fabric-config" "$FABLO_NETWORK_ROOT/fabric-config/config" "Group1Genesis"

  # Create directory for chaincode packages to avoid permission errors on linux
  mkdir -p "$FABLO_NETWORK_ROOT/fabric-config/chaincode-packages"
}

startNetwork() {
  printHeadline "Starting network" "U1F680"
  (cd "$FABLO_NETWORK_ROOT"/fabric-docker && docker-compose up -d)
  sleep 4
}

generateChannelsArtifacts() {
  printHeadline "Generating config for 'my-channel'" "U1F913"
  createChannelTx "my-channel" "$FABLO_NETWORK_ROOT/fabric-config" "MyChannel" "$FABLO_NETWORK_ROOT/fabric-config/config"
}

installChannels() {
  printHeadline "Creating 'my-channel' on CountryA/peer0" "U1F63B"
  docker exec -i cli.countrya.example.com bash -c "source scripts/channel_fns.sh; createChannelAndJoinTls 'my-channel' 'CountryAMSP' 'peer0.countrya.example.com:7041' 'crypto/users/Admin@countrya.example.com/msp' 'crypto/users/Admin@countrya.example.com/tls' 'crypto-orderer/tlsca.orderer.example.com-cert.pem' 'orderer0.group1.orderer.example.com:7030';"

  printItalics "Joining 'my-channel' on  CountryB/peer0" "U1F638"
  docker exec -i cli.countryb.example.com bash -c "source scripts/channel_fns.sh; fetchChannelAndJoinTls 'my-channel' 'CountryBMSP' 'peer0.countryb.example.com:7061' 'crypto/users/Admin@countryb.example.com/msp' 'crypto/users/Admin@countryb.example.com/tls' 'crypto-orderer/tlsca.orderer.example.com-cert.pem' 'orderer0.group1.orderer.example.com:7030';"
}

installChaincodes() {
  if [ -n "$(ls "$CHAINCODES_BASE_DIR/./chaincode-go")" ]; then
    local version="0.0.1"
    printHeadline "Packaging chaincode 'chaincode1'" "U1F60E"
    chaincodeBuild "chaincode1" "golang" "$CHAINCODES_BASE_DIR/./chaincode-go" "16"
    chaincodePackage "cli.countrya.example.com" "peer0.countrya.example.com:7041" "chaincode1" "$version" "golang" printHeadline "Installing 'chaincode1' for CountryA" "U1F60E"
    chaincodeInstall "cli.countrya.example.com" "peer0.countrya.example.com:7041" "chaincode1" "$version" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
    chaincodeApprove "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" ""
    printHeadline "Installing 'chaincode1' for CountryB" "U1F60E"
    chaincodeInstall "cli.countryb.example.com" "peer0.countryb.example.com:7061" "chaincode1" "$version" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
    chaincodeApprove "cli.countryb.example.com" "peer0.countryb.example.com:7061" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" ""
    printItalics "Committing chaincode 'chaincode1' on channel 'my-channel' as 'CountryA'" "U1F618"
    chaincodeCommit "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" "peer0.countrya.example.com:7041,peer0.countryb.example.com:7061" "crypto-peer/peer0.countrya.example.com/tls/ca.crt,crypto-peer/peer0.countryb.example.com/tls/ca.crt" ""
  else
    echo "Warning! Skipping chaincode 'chaincode1' installation. Chaincode directory is empty."
    echo "Looked in dir: '$CHAINCODES_BASE_DIR/./chaincode-go'"
  fi

}

installChaincode() {
  local chaincodeName="$1"
  if [ -z "$chaincodeName" ]; then
    echo "Error: chaincode name is not provided"
    exit 1
  fi

  local version="$2"
  if [ -z "$version" ]; then
    echo "Error: chaincode version is not provided"
    exit 1
  fi

  if [ "$chaincodeName" = "chaincode1" ]; then
    if [ -n "$(ls "$CHAINCODES_BASE_DIR/./chaincode-go")" ]; then
      printHeadline "Packaging chaincode 'chaincode1'" "U1F60E"
      chaincodeBuild "chaincode1" "golang" "$CHAINCODES_BASE_DIR/./chaincode-go" "16"
      chaincodePackage "cli.countrya.example.com" "peer0.countrya.example.com:7041" "chaincode1" "$version" "golang" printHeadline "Installing 'chaincode1' for CountryA" "U1F60E"
      chaincodeInstall "cli.countrya.example.com" "peer0.countrya.example.com:7041" "chaincode1" "$version" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
      chaincodeApprove "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" ""
      printHeadline "Installing 'chaincode1' for CountryB" "U1F60E"
      chaincodeInstall "cli.countryb.example.com" "peer0.countryb.example.com:7061" "chaincode1" "$version" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
      chaincodeApprove "cli.countryb.example.com" "peer0.countryb.example.com:7061" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" ""
      printItalics "Committing chaincode 'chaincode1' on channel 'my-channel' as 'CountryA'" "U1F618"
      chaincodeCommit "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" "peer0.countrya.example.com:7041,peer0.countryb.example.com:7061" "crypto-peer/peer0.countrya.example.com/tls/ca.crt,crypto-peer/peer0.countryb.example.com/tls/ca.crt" ""

    else
      echo "Warning! Skipping chaincode 'chaincode1' install. Chaincode directory is empty."
      echo "Looked in dir: '$CHAINCODES_BASE_DIR/./chaincode-go'"
    fi
  fi
}

runDevModeChaincode() {
  local chaincodeName=$1
  if [ -z "$chaincodeName" ]; then
    echo "Error: chaincode name is not provided"
    exit 1
  fi

  if [ "$chaincodeName" = "chaincode1" ]; then
    local version="0.0.1"
    printHeadline "Approving 'chaincode1' for CountryA (dev mode)" "U1F60E"
    chaincodeApprove "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "0.0.1" "orderer0.group1.orderer.example.com:7030" "" "false" "" ""
    printHeadline "Approving 'chaincode1' for CountryB (dev mode)" "U1F60E"
    chaincodeApprove "cli.countryb.example.com" "peer0.countryb.example.com:7061" "my-channel" "chaincode1" "0.0.1" "orderer0.group1.orderer.example.com:7030" "" "false" "" ""
    printItalics "Committing chaincode 'chaincode1' on channel 'my-channel' as 'CountryA' (dev mode)" "U1F618"
    chaincodeCommit "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "0.0.1" "orderer0.group1.orderer.example.com:7030" "" "false" "" "peer0.countrya.example.com:7041,peer0.countryb.example.com:7061" "" ""

  fi
}

upgradeChaincode() {
  local chaincodeName="$1"
  if [ -z "$chaincodeName" ]; then
    echo "Error: chaincode name is not provided"
    exit 1
  fi

  local version="$2"
  if [ -z "$version" ]; then
    echo "Error: chaincode version is not provided"
    exit 1
  fi

  if [ "$chaincodeName" = "chaincode1" ]; then
    if [ -n "$(ls "$CHAINCODES_BASE_DIR/./chaincode-go")" ]; then
      printHeadline "Packaging chaincode 'chaincode1'" "U1F60E"
      chaincodeBuild "chaincode1" "golang" "$CHAINCODES_BASE_DIR/./chaincode-go" "16"
      chaincodePackage "cli.countrya.example.com" "peer0.countrya.example.com:7041" "chaincode1" "$version" "golang" printHeadline "Installing 'chaincode1' for CountryA" "U1F60E"
      chaincodeInstall "cli.countrya.example.com" "peer0.countrya.example.com:7041" "chaincode1" "$version" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
      chaincodeApprove "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" ""
      printHeadline "Installing 'chaincode1' for CountryB" "U1F60E"
      chaincodeInstall "cli.countryb.example.com" "peer0.countryb.example.com:7061" "chaincode1" "$version" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
      chaincodeApprove "cli.countryb.example.com" "peer0.countryb.example.com:7061" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" ""
      printItalics "Committing chaincode 'chaincode1' on channel 'my-channel' as 'CountryA'" "U1F618"
      chaincodeCommit "cli.countrya.example.com" "peer0.countrya.example.com:7041" "my-channel" "chaincode1" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "crypto-orderer/tlsca.orderer.example.com-cert.pem" "peer0.countrya.example.com:7041,peer0.countryb.example.com:7061" "crypto-peer/peer0.countrya.example.com/tls/ca.crt,crypto-peer/peer0.countryb.example.com/tls/ca.crt" ""

    else
      echo "Warning! Skipping chaincode 'chaincode1' upgrade. Chaincode directory is empty."
      echo "Looked in dir: '$CHAINCODES_BASE_DIR/./chaincode-go'"
    fi
  fi
}

notifyOrgsAboutChannels() {
  printHeadline "Creating new channel config blocks" "U1F537"
  createNewChannelUpdateTx "my-channel" "CountryAMSP" "MyChannel" "$FABLO_NETWORK_ROOT/fabric-config" "$FABLO_NETWORK_ROOT/fabric-config/config"
  createNewChannelUpdateTx "my-channel" "CountryBMSP" "MyChannel" "$FABLO_NETWORK_ROOT/fabric-config" "$FABLO_NETWORK_ROOT/fabric-config/config"

  printHeadline "Notyfing orgs about channels" "U1F4E2"
  notifyOrgAboutNewChannelTls "my-channel" "CountryAMSP" "cli.countrya.example.com" "peer0.countrya.example.com" "orderer0.group1.orderer.example.com:7030" "crypto-orderer/tlsca.orderer.example.com-cert.pem"
  notifyOrgAboutNewChannelTls "my-channel" "CountryBMSP" "cli.countryb.example.com" "peer0.countryb.example.com" "orderer0.group1.orderer.example.com:7030" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  printHeadline "Deleting new channel config blocks" "U1F52A"
  deleteNewChannelUpdateTx "my-channel" "CountryAMSP" "cli.countrya.example.com"
  deleteNewChannelUpdateTx "my-channel" "CountryBMSP" "cli.countryb.example.com"
}

printStartSuccessInfo() {
  printHeadline "Done! Enjoy your fresh network" "U1F984"
}

stopNetwork() {
  printHeadline "Stopping network" "U1F68F"
  (cd "$FABLO_NETWORK_ROOT"/fabric-docker && docker-compose stop)
  sleep 4
}

networkDown() {
  printHeadline "Destroying network" "U1F916"
  (cd "$FABLO_NETWORK_ROOT"/fabric-docker && docker-compose down)

  printf "\nRemoving chaincode containers & images... \U1F5D1 \n"
  for container in $(docker ps -a | grep "dev-peer0.countrya.example.com-chaincode1" | awk '{print $1}'); do
    echo "Removing container $container..."
    docker rm -f "$container" || echo "docker rm of $container failed. Check if all fabric dockers properly was deleted"
  done
  for image in $(docker images "dev-peer0.countrya.example.com-chaincode1*" -q); do
    echo "Removing image $image..."
    docker rmi "$image" || echo "docker rmi of $image failed. Check if all fabric dockers properly was deleted"
  done
  for container in $(docker ps -a | grep "dev-peer0.countryb.example.com-chaincode1" | awk '{print $1}'); do
    echo "Removing container $container..."
    docker rm -f "$container" || echo "docker rm of $container failed. Check if all fabric dockers properly was deleted"
  done
  for image in $(docker images "dev-peer0.countryb.example.com-chaincode1*" -q); do
    echo "Removing image $image..."
    docker rmi "$image" || echo "docker rmi of $image failed. Check if all fabric dockers properly was deleted"
  done

  printf "\nRemoving generated configs... \U1F5D1 \n"
  rm -rf "$FABLO_NETWORK_ROOT/fabric-config/config"
  rm -rf "$FABLO_NETWORK_ROOT/fabric-config/crypto-config"
  rm -rf "$FABLO_NETWORK_ROOT/fabric-config/chaincode-packages"

  printHeadline "Done! Network was purged" "U1F5D1"
}
