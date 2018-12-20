#!/bin/sh

# setup darknode, it expects there already has an config file in the folder.
mv ./darknode $HOME/.darknode/bin/darknode
systemctl --user restart darknode.service
