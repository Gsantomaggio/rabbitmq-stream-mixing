rm -rf rabbitmq-stream-go-client
git clone https://github.com/rabbitmq/rabbitmq-stream-go-client.git
docker build -f rabbitmq-stream-go-client/examples/reliable/Dockerfile -t stream-reliable-client rabbitmq-stream-go-client
rm -rf rabbitmq-stream-go-client
docker tag stream-reliable-client:latest localhost:5000/stream-reliable-client:latest
docker push localhost:5000/stream-reliable-client:latest
