#!/bin/bash
# how to run this script:
#   cd confd && sh integration/json/test.sh
#   then can see: cat /tmp/confd-*.conf

jsonname="./.__confd-test.json"
cat > $jsonname << DELIM
{
  "prefix": "prefix",
  "keys": [
    {"key": "database/host", "value": "127.0.0.1"},
    {"key": "database/password", "value": "p@sSw0rd"},
    {"key": "database/port", "value": "3306"},
    {"key": "database/username", "value": "confd"},
    {"key": "upstream/app1", "value": "10.0.1.10:8080"},
    {"key": "upstream/app2", "value": "10.0.1.11:8080"}
  ],
  "fullkeys": [
    {"fullkey": "database/host", "value": "127.0.0.1"},
    {"fullkey": "database/password", "value": "p@sSw0rd"},
    {"fullkey": "database/port", "value": "3306"},
    {"fullkey": "database/username", "value": "confd"},
    {"fullkey": "upstream/app1", "value": "10.0.1.10:8080"},
    {"fullkey": "upstream/app2", "value": "10.0.1.11:8080"}
  ],
  "fasdfsilasfdalkj": [
    {"fullkey": "key", "value": "foobar"}
  ]
}
DELIM


# Run confd
confd -onetime -verbose -debug -confdir ./integration/confdir -backend json -node $jsonname
confd -onetime -verbose -debug -confdir ./integration/confdir -backend json -node $jsonname -prefix prefix

rm $jsonname
