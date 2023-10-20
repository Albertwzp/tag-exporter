OPT := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BINARY := tag-exporter

.PHONY: cross
cross: $(BINARY)

.PHONY: clean
clean:
	@rm -rf $(BINARY)

$(BINARY):
	${OPT} go build -o $(BINARY) -ldflags="-s -w" ./main.go