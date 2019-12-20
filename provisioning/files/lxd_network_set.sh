#!/bin/bash

key=$1
value=$2

printf "$value" | lxc network set lxdbr0 $1 -
