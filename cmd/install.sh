#!/bin/sh

generate_updater_service() {
  echo "[Unit]
Description=Republic Protocol's Darknode Automatic Updater
AssertPathExists=$HOME/.darknode

[Service]
ExecStart=/bin/bash $HOME/.darknode/bin/auto-updater.sh
Restart=on-failure
PrivateTmp=true
NoNewPrivileges=true

# Specifies which signal to use when killing a service. Defaults to SIGTERM.
# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)
KillSignal=SIGHUP

[Install]
WantedBy=default.target" > $HOME/.config/systemd/user/darknode-updater.service
}

generate_darknode_service() {
  echo "[Unit]
Description=Republic Protocol's Darknode Daemon
AssertPathExists=$HOME/.darknode

[Service]
WorkingDirectory=$HOME/.darknode
ExecStart=$HOME/.darknode/bin/darknode --config $HOME/.darknode/config.json
Restart=on-failure
PrivateTmp=true
NoNewPrivileges=true

# Specifies which signal to use when killing a service. Defaults to SIGTERM.
# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)
KillSignal=SIGHUP

[Install]
WantedBy=default.target" > $HOME/.config/systemd/user/darknode.service
}

# generate service files
mkdir -p $HOME/.config/systemd/user
generate_updater_service
generate_darknode_service

# setup darknode, it expects there already has an config file in the folder.
mkdir -p $HOME/.darknode/bin
mv ./darknode $HOME/.darknode/bin/darknode
mv ./auto-updater.sh $HOME/.darknode/bin/auto-updater.sh


# Start services
loginctl enable-linger `whoami`
systemctl --user enable darknode.service
systemctl --user enable darknode-updater.service
systemctl --user start darknode.service
systemctl --user start darknode-updater.service
