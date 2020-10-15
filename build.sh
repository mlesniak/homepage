#!/usr/bin/env sh

if [  -e aperol ]
then
    aperol
else
    go run main.go
fi

rsync -rv docs/ server:/root/www