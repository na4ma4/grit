GO_MATRIX_OS := darwin linux windows

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/grit
	GRIT_CONFIG=etc/testing.toml "$<" $(RUN_ARGS)

homebrew: artifacts/archives/grit-$(GIT_HEAD_TAG)-darwin-amd64.zip
	bin/homebrew.sh "$(GIT_HEAD_TAG)" "$<"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
