.PHONY: run
.EXPORT_ALL_VARIABLES:

MAINFILES = ./...

run:
	@go run $(MAINFILES)
