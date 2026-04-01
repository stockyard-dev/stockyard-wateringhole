.PHONY: build run clean
build:
	CGO_ENABLED=0 go build -o wateringhole ./cmd/wateringhole/
run: build
	./wateringhole
clean:
	rm -f wateringhole
