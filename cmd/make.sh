mkdir -p ../releases/darknode

cd darknode
xgo --targets=linux/amd64 .
mv darknode-linux-amd64 ../../releases/darknode/darknode

cd ../
cp install.sh ../releases/darknode
cp auto-updater.sh ../releases/darknode
cp update.sh ../releases/darknode

cd ../releases
zip -r darknode-linux-amd64.zip ./darknode