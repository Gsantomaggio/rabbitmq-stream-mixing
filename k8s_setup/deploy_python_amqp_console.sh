rm -rf rabbitmq-amqp-python-client
git clone --branch v2_main https://github.com/rabbitmq/rabbitmq-amqp-python-client.git
docker build -f rabbitmq-amqp-python-client/docs/examples/docker/Dockerfile -t python-amqp-reliable-client rabbitmq-amqp-python-client
rm -rf rabbitmq-amqp-python-client
docker tag python-amqp-console-client:latest localhost:5000/python-amqp-reliable-client:latest
docker push localhost:5000/python-amqp-reliable-client:latest
