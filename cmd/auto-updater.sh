#!/bin/sh

get_latest_release() {
  curl -s https://api.github.com/repos/republicprotocol/darknode-cli/releases/latest \
    | grep "browser_download_url.*darknode-$1.zip" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi -
  mv darknode-$1.zip darknode.zip
}

while true
do
  R=$(($RANDOM%60))
  if test $R -eq 0; then
    echo "Updating darknode..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S)

    get_latest_release linux-amd64
    unzip darknode.zip
    cd darknode
    chmod +x update.sh
    ./update.sh
    cd
    rm -rf darknode
    rm darknode.zip

    echo $timestamp >> /home/darknode/.darknode/update.log
    echo "Finish updating"
  fi
  sleep 1m
done