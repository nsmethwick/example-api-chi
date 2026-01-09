.DEFAULT_GOAL := run        

.PHONY: help
help: ## Displays list of Makefile targets and documented variables
	@echo "Targets:"
	@MAX_LEN_TARGET=$$(grep -h -E '^[0-9a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "; max=0} {len=length($$1); if (len>max) max=len} END {print max}'); \
	grep -h -E '^[0-9a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk -v max_len="$$MAX_LEN_TARGET" 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-" max_len "s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Variables:"
	@MAX_LEN_VAR=$$(grep -h -E '^[0-9a-zA-Z_.-]+\s[?:]?=.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "[?:]?=.*?## "; max=0} {len=length($$1); if (len>max) max=len} END {print max}'); \
	grep -h -E '^[0-9a-zA-Z_.-]+\s[?:]?=.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk -v max_len="$$MAX_LEN_VAR" 'BEGIN {FS = "[?:]?=.*?## "}; {printf "  \033[36m%-" max_len "s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Default target:"
	@printf "  \033[36m%s\033[0m\n" $(.DEFAULT_GOAL)

.PHONY: run
run: ## Run the API server
	go run ./cmd/api/main.go