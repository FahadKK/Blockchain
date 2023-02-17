#!/usr/bin/env bash

source "$FABLO_NETWORK_ROOT/fabric-docker/scripts/channel-query-functions.sh"

set -eu

channelQuery() {
  echo "-> Channel query: " + "$@"

  if [ "$#" -eq 1 ]; then
    printChannelsHelp

  elif [ "$1" = "list" ] && [ "$2" = "countrya" ] && [ "$3" = "peer0" ]; then

    peerChannelListTls "cli.countrya.example.com" "peer0.countrya.example.com:7041" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  elif
    [ "$1" = "list" ] && [ "$2" = "countryb" ] && [ "$3" = "peer0" ]
  then

    peerChannelListTls "cli.countryb.example.com" "peer0.countryb.example.com:7061" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  elif

    [ "$1" = "getinfo" ] && [ "$2" = "my-channel" ] && [ "$3" = "countrya" ] && [ "$4" = "peer0" ]
  then

    peerChannelGetInfoTls "my-channel" "cli.countrya.example.com" "peer0.countrya.example.com:7041" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  elif [ "$1" = "fetch" ] && [ "$2" = "config" ] && [ "$3" = "my-channel" ] && [ "$4" = "countrya" ] && [ "$5" = "peer0" ]; then
    TARGET_FILE=${6:-"$channel-config.json"}

    peerChannelFetchConfigTls "my-channel" "cli.countrya.example.com" "$TARGET_FILE" "peer0.countrya.example.com:7041" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  elif [ "$1" = "fetch" ] && [ "$3" = "my-channel" ] && [ "$4" = "countrya" ] && [ "$5" = "peer0" ]; then
    BLOCK_NAME=$2
    TARGET_FILE=${6:-"$BLOCK_NAME.block"}

    peerChannelFetchBlockTls "my-channel" "cli.countrya.example.com" "${BLOCK_NAME}" "peer0.countrya.example.com:7041" "crypto-orderer/tlsca.orderer.example.com-cert.pem" "$TARGET_FILE"

  elif
    [ "$1" = "getinfo" ] && [ "$2" = "my-channel" ] && [ "$3" = "countryb" ] && [ "$4" = "peer0" ]
  then

    peerChannelGetInfoTls "my-channel" "cli.countryb.example.com" "peer0.countryb.example.com:7061" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  elif [ "$1" = "fetch" ] && [ "$2" = "config" ] && [ "$3" = "my-channel" ] && [ "$4" = "countryb" ] && [ "$5" = "peer0" ]; then
    TARGET_FILE=${6:-"$channel-config.json"}

    peerChannelFetchConfigTls "my-channel" "cli.countryb.example.com" "$TARGET_FILE" "peer0.countryb.example.com:7061" "crypto-orderer/tlsca.orderer.example.com-cert.pem"

  elif [ "$1" = "fetch" ] && [ "$3" = "my-channel" ] && [ "$4" = "countryb" ] && [ "$5" = "peer0" ]; then
    BLOCK_NAME=$2
    TARGET_FILE=${6:-"$BLOCK_NAME.block"}

    peerChannelFetchBlockTls "my-channel" "cli.countryb.example.com" "${BLOCK_NAME}" "peer0.countryb.example.com:7061" "crypto-orderer/tlsca.orderer.example.com-cert.pem" "$TARGET_FILE"

  else

    echo "$@"
    echo "$1, $2, $3, $4, $5, $6, $7, $#"
    printChannelsHelp
  fi

}

printChannelsHelp() {
  echo "Channel management commands:"
  echo ""

  echo "fablo channel list countrya peer0"
  echo -e "\t List channels on 'peer0' of 'CountryA'".
  echo ""

  echo "fablo channel list countryb peer0"
  echo -e "\t List channels on 'peer0' of 'CountryB'".
  echo ""

  echo "fablo channel getinfo my-channel countrya peer0"
  echo -e "\t Get channel info on 'peer0' of 'CountryA'".
  echo ""
  echo "fablo channel fetch config my-channel countrya peer0 [file-name.json]"
  echo -e "\t Download latest config block and save it. Uses first peer 'peer0' of 'CountryA'".
  echo ""
  echo "fablo channel fetch <newest|oldest|block-number> my-channel countrya peer0 [file name]"
  echo -e "\t Fetch a block with given number and save it. Uses first peer 'peer0' of 'CountryA'".
  echo ""

  echo "fablo channel getinfo my-channel countryb peer0"
  echo -e "\t Get channel info on 'peer0' of 'CountryB'".
  echo ""
  echo "fablo channel fetch config my-channel countryb peer0 [file-name.json]"
  echo -e "\t Download latest config block and save it. Uses first peer 'peer0' of 'CountryB'".
  echo ""
  echo "fablo channel fetch <newest|oldest|block-number> my-channel countryb peer0 [file name]"
  echo -e "\t Fetch a block with given number and save it. Uses first peer 'peer0' of 'CountryB'".
  echo ""

}
