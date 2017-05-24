.PHONY:	install
install:
	go get -u github.com/kardianos/govendor
	govendor sync

.PHONY:	errors
errors:
	go get -u github.com/kisielk/errcheck
	errcheck ./cmd/...
	errcheck ./pkg/...

.PHONY:	verify
verify: errors