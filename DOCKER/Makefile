build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build_testing:
	docker build --tag reapchain-core/testing -f ./Dockerfile.testing .

build_amazonlinux_buildimage:
	docker build -t "reapchain/reapchain-core:build_c-amazonlinux" -f Dockerfile.build_c-amazonlinux .

.PHONY: build push build_testing build_amazonlinux_buildimage
