.PHONY : ci-test
ci-test:
	make -C ./svt-app-operator/ ci-test


.PHONY : ci-install
ci-install:
	git --version
	go version
	kubectl version
	docker version
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep version
	echo "$${GOPATH}"
	install_operator_sdk.sh


.PHONY : ci-script
ci-script: ci-install ci-test
