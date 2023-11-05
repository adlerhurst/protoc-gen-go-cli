generate-option:
	$(RM) -r gen/proto/adlerhurst
	buf generate --path proto/adlerhurst

.PHONY: compile
compile:
	go install .

generate-example: compile
	$(RM) -r gen/proto/example
	buf generate --path proto/example
