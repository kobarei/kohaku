all:
	go build -o bin/kohaku cmd/kohaku/main.go

run: clean all
	./bin/kohaku -c config.toml

clean:
	rm -f bin/kohaku
	rm -rf ./kohaku.log

test:
	go test -race -v
