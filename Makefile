.PHONY : ci-test
ci-test:
	make -C ./svt-app-operator/ ci-test


.PHONY : ci-install
ci-install:
	git --version
	go version
	echo $$PATH
	./script/ci/install_dep.sh
	kubectl version
	docker version
	#curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep version
	echo "$${GOPATH}"
	./script/ci/install_dep.sh


.PHONY : ci-script
ci-script: ci-install ci-test
