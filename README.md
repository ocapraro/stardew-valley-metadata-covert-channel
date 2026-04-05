# Stardew Valley Metadata Covert Channel
Using the inventory data of players in Stardew Valley to communicate a covert message.

## Run the decoder
```zsh
cd servers/src/decoder
set -a
source ../.env
set +a
go run main.go
```

## Run the encoder
```zsh
cd servers/src/encoder
go run main.go
```