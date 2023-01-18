.PHONY: producer consumer all



producer:
	@echo "running python-rstream producer"
	venv/bin/python3 python-stream-client/producer.py

consumer:
	@echo "running python-rstream consumer"
	venv/bin/python3 python-stream-client/consumer.py


init-python:
	@echo "init python venv"
	rm -rf venv
	python3 -m venv venv
	venv/bin/pip install -r requirements.txt

	
	
all: producer consumer