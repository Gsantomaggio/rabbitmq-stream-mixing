rm -rf rabbitmq-stream-dotnet-client
git clone https://github.com/rabbitmq/rabbitmq-stream-dotnet-client.git
docker build -f rabbitmq-stream-dotnet-client/docs/ReliableClient/Dockerfile -t dotnet-stream-reliable-client rabbitmq-stream-dotnet-client
rm -rf rabbitmq-stream-dotnet-client
docker tag dotnet-stream-reliable-client:latest localhost:5000/dotnet-stream-reliable-client:latest
docker push localhost:5000/dotnet-stream-reliable-client:latest
