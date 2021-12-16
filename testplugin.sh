#! /usr/bin/env bash
if [ -z $I2P ]; then
	I2P=$HOME/i2p
fi
if [ -z $CONFIG ]; then
	CONFIG=$HOME/.i2p
fi
if [ -z $1 ]; then
	echo "Please enter the plugin name(example: railroad-linux)"
	exit 1
fi
export PLUGINNAME="$1"
export PLUGIN="$CONFIG/plugins/$PLUGINNAME"
export PLUGINPATH="$PLUGIN/lib/$PLUGINNAME"
export PREARGS="$(cat "$CONFIG/plugins/$PLUGINNAME/clients.config" | grep 'clientApp.0.args' | sed 's|clientApp.0.args=||g' | cut -sf6- -d'-')"
export ARGS=$(echo $PREARGS | sed "s|\$PLUGIN|$PLUGIN|g")
echo "Testing: $PLUGINNAME"
echo "Using environment:"
echo "  I2P=$I2P"
echo "  CONFIG=$CONFIG"
echo "  PLUGIN=$PLUGIN"
echo "  PLUGINPATH=$PLUGINPATH"
echo "  ARGS=$ARGS"
echo "  COMMAND="$PLUGINPATH" -$ARGS"
"$PLUGINPATH" -$ARGS > log.log 2> err.log &
tail -f log.log err.log

