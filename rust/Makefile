.PHONY: producer consumer all


producer:
	@echo "running rust producer"
	cargo run  --manifest-path producer/Cargo.toml

consumer:
	@echo "running rust consumer"
	cargo run  --manifest-path consumer/Cargo.toml

all: producer consumer
