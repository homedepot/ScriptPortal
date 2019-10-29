#!/usr/bin/env bash
##put the script portal in the bin folder
install build/ScriptPortal /usr/bin
#move main config to var/lib because that's what man heir says
mkdir -p /var/lib/ScriptPortal
if [ ! -f /var/lib/ScriptPortal/scriptconfig.json ]; then
	cp config/scriptconfig.json /var/lib/ScriptPortal/scriptConfig.json
fi
#put the templates in /usr/share because that what man heir says
mkdir -p /usr/share/ScriptPortal
cp -r templates /usr/share/ScriptPortal
#put the compiled plugins in /var/opt because that's what man heir says
if [ -f build/*.so ]
then
	mkdir -p /var/opt/ScriptPortal;	
	cp build/*.so /var/opt/ScriptPortal/
fi
