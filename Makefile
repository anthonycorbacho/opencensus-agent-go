PROTO_FILES = $(wildcard api/proto/*.proto)

proto	:
	@for proto in $(PROTO_FILES); do \
	  if ! protoc --proto_path=api/proto --go_out=plugins=grpc:pkg/api "$$proto" ; then \
            echo " ❎️  Couldnt generate go file for $$proto️"; \
            exit 1; \
          else \
            echo " ✅  $$proto"; \
          fi \
	done
