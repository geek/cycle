PKG := github.com/geek/cycle

fmt:
	@test -z "$(shell gofmt -s -l -d -e ./ | tee /dev/stderr)"

vet:
	go vet $(PKG)

test: vet fmt
	go test $(PKG)
