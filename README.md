# Overview

Data and build program to generate content for [mlesniak.com](https://mlesniak.com).

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