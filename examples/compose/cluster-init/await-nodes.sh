#!/usr/bin/env sh

available_nodes=0

./wait-for 172.16.238.11:5984 -t 10
result=$?
if [ ${result} -eq 0 ] ; then
  available_nodes=$((available_nodes+1))
fi

./wait-for 172.16.238.12:5984 -t 10
result=$?
if [ ${result} -eq 0 ] ; then
  available_nodes=$((available_nodes+1))
fi

./wait-for 172.16.238.13:5984 -t 10
result=$?
if [ ${result} -eq 0 ] ; then
  available_nodes=$((available_nodes+1))
fi

if [ ${available_nodes} -eq 3 ]
then echo "alright"
else echo "missing node, got $available_nodes nodes"
fi
