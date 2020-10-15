# Redis is not reachable with a custom configuration

#lessonslearned

If you use the [standard docker image for redis](https://hub.docker.com/_/redis) and switch to a custom configuration file, 
i.e. start the image with

    docker run -it --rm --name redis -p 6379:6379 \
        -v $(pwd)/data:/data \
        -v $(pwd)/redis.conf:/usr/local/etc/redis/redis.conf \
        redis:6 \
        redis-server /usr/local/etc/redis/redis.conf --appendonly yes

it solely listens to the loopback device. Hence, you explicitly have to disable this
by commenting the `bind` setting in `redis.conf`, i.e. find the line `bind 127.0.0.1` and write

    # bind 127.0.0.1

If you do this, please activate [ACLs](https://redis.io/topics/acl), otherwise everyone can connect to your database.