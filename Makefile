PKG=github.com/boynton/midi-ell
BIN=$(GOPATH)/bin/mell

all:
	CGO_CFLAGS="-I/opt/homebrew/include" CGO_LDFLAGS="-L/opt/homebrew/lib" go install $(PKG)/cmd/mell

prebuild:
	go get -u github.com/boynton/ell
	go get -u github.com/boynton/repl
#   brew install portmidi
	CGO_CFLAGS="-I/opt/homebrew/include" CGO_LDFLAGS="-L/opt/homebrew/lib" go get -u github.com/rakyll/portmidi

test:
	go run miditest.go

run:
	./bin/mell


clean:
	go clean $(PKG)
	rm -rf *~ $(BIN)

check:
	@(cd $(GOPATH)/src/$(PKG); go vet $(PKG); go fmt $(PKG))
