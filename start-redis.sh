#docker network create blocopadnet
docker run --rm \
    -p 6379:6379 \
    -v ./data:/data \
    --name redisbase \
    --network blocopadnet \
    redis