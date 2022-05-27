testcase:
	@echo "Updating testcase git submodule..."
	git submodule update --init --recursive --remote
	make -C testcase-compiler testcase.out

check: testcase
	go test ./...

verify-commits: testcase
	bash ci/verify_commits.sh
