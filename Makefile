.PHONY : ci-test
ci-test:
	make -C ./svt-app-operator/ ci-test


.PHONY : ci-install
ci-install:
	git --version
	go version
	#snap --version
	echo $${PATH}
	./script/ci/install_dep.sh
	kubectl version --client=true
	docker version
	dep version
	echo "$${GOPATH}"



.PHONY : ci-script
ci-script: ci-install ci-test
