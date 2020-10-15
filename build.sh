#!/usr/bin/env sh

if [ ! -e aperol ]
then
    go build
fi
aperol
git add .
git commit -m "$(date)"
git push


rsync -rv docs/ server:/root/www