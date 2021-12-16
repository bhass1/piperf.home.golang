#!/bin/bash
iperf3 -c iperf.he.net -n 3M -N -R -J -p 5201
if [ $? -eq 1 ]; then
	iperf3 -c iperf.scottlinux.com -n 3M -N -R -J -p 5201
else
	exit 0
fi
if [ $? -eq 1 ]; then
	iperf3 -c ping.online.net -n 3M -N -R -J -p 5207
else
	exit 0
fi
if [ $? -eq 1 ]; then
	iperf3 -c bouygues.iperf.fr -n 3M -N -R -J -p 5207
else
	exit 0
fi
