docker run --rm \
    -p 8080:8080 \
    --name blocopad \
    --network blocopadnet \
    -v ./keys:/keys \
    --env API_DB_URL=redisbase:6379 \
    --env API_PRIVATE_KEY=/keys/key \
    --env API_PUBLIC_KEY=/keys/key.pub \
    blocopad:v001
