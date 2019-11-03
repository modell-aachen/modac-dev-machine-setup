#!/bin/bash

key=$1
value=$2

echo -e "$value" | lxc network set lxdbr0 $1 -
