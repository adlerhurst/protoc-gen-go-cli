generate-option:
	$(RM) -r gen
	buf generate --path proto/adlerhurst

.PHONY: compile
compile:
	go install .

generate-example: compile
	$(RM) -r proto/example/api
	buf generate --path proto/example
