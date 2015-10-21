PKG=github.com/boynton/midi-ell
BIN=$(GOPATH)/bin/mell

all:
	go install $(PKG)/cmd/mell

prebuild:
	go get -u github.com/boynton/ell
	go get -u github.com/boynton/repl
	go get -u github.com/pborman/uuid
	go get -u github.com/rakyll/portmidi

test:
	go run miditest.go

run:
	./bin/mell


clean:
	go clean $(PKG)
	rm -rf *~ $(BIN)

check:
	@(cd $(GOPATH)/src/$(PKG); go vet $(PKG); go fmt $(PKG))
