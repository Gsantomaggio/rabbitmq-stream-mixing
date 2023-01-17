.PHONY: producer consumer all


producer:
	@echo "running java producer"
	mvn install -f ./producer/pom.xml
	java -jar ./producer/target/rabbitmq-stream-client-1.0-SNAPSHOT-jar-with-dependencies.jar

consumer:
	@echo "running java consumer"
	mvn install -f ./consumer/pom.xml
	java -jar ./consumer/target/rabbitmq-stream-client-1.0-SNAPSHOT-jar-with-dependencies.jar

all: producer consumer