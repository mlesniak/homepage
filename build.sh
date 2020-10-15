#!/usr/bin/env sh

git add .
git commit -m "$(date)"
git push

if [ -e aperol ]
then
    aperol
else
    go run main.go
fi

rsync -rv docs/ server:/root/www