.PHONY: proto tidy

## Gera código Go a partir dos .proto (requer buf: https://buf.build/docs/installation)
proto:
	buf generate

## Baixa dependências e atualiza go.sum
tidy:
	go mod tidy
