testcase:
	@echo "Updating testcase git submodule..."
	git submodule update --remote
	make -C testcase-compiler testcase.out

check: testcase
	go test ./...
