.PHONY : ci-test
ci-test:
	make -C ./svt-app-operator/ ci-test


.PHONY : ci-install
ci-install:
	git --version
	go version
	python -V
	pip --version
	python3 -V
	pip3 --version
	docker version
	echo $${PATH}
	echo "$${GOPATH}"
	./script/ci/install_dep.sh
	kubectl version --client=true
	dep version





.PHONY : ci-script
ci-script: ci-install ci-test
