docker run --rm \
    -p 6379:6379 \
    -v ./data:/data \
    --name redis-api-go \
    redis