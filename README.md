# Overview

Quick, hacky and dirty build program to generate content for [mlesniak.com](https://mlesniak.com). This code is highly specialized to my needs; I've probably written much more blog engines than I ever will write blog posts...

## webserver.sh

    docker stop -t 0 webserver
    docker rm webserver
    cp /etc/letsencrypt/live/mlesniak.com-0001/privkey.pem certificates/
    cp /etc/letsencrypt/live/mlesniak.com-0001/fullchain.pem certificates/
    docker run \
            --name=webserver \
            -d \
            --volume="/root/certificates/:/etc/ssl/local/" \
            --volume="/root/www/:/var/www/" \
            --volume=$(pwd)/default.conf:/etc/nginx/conf.d/default.conf \
            -p 8080:8080 -p 80:80 -p 443:443 \
            nginx:latest

## default.conf

    server {
    	listen 80       default_server;
    	listen [::]:80  default_server;
    	server_name _;
    	return 301 https://$host$request_uri;
    }
    
    server {
        listen          443;
        server_name     mlesniak.com www.mlesniak.com;
        root            /var/www;
        ssl                 on;
        ssl_certificate     /etc/ssl/local/fullchain.pem;
        ssl_certificate_key /etc/ssl/local/privkey.pem;
    }

## tasks.json

    {
        "version": "2.0.0",
        "tasks": [
            {
                "label": "publish",
                "type": "shell",
                "options": {
                    "cwd": ".."
                },
                "presentation": {
                    "echo": true,
                    "reveal": "never",
                    "focus": false,
                    "panel": "shared",
                    "showReuseMessage": true,
                    "clear": true
                },
                "command": "build.sh",
                "problemMatcher": [],
                "group": {
                    "kind": "build",
                    "isDefault": true
                }
            }
        ]
    }
