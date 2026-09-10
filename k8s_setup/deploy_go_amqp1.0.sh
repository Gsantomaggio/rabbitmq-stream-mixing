git clone https://github.com/rabbitmq/rabbitmq-amqp-go-client.git .
cd rabbitmq-amqp-go-client && make docker-build-reliable-example
rm -rf rabbitmq-amqp-go-client
docker tag go-amqp1.0-client:latest localhost:5000/go-amqp1.0-client:latest
docker push localhost:5000/go-amqp1.0-client:latest