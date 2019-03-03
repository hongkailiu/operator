.PHONY : ci-test
ci-test:
	echo "ci-test for svt-app-operator"
	kubectl create namespace operator-test
	operator-sdk build docker.io/hongkailiu/svt-app-operator:ci
	docker images
	operator-sdk test local ./test/e2e --namespace operator-test --image docker.io/hongkailiu/svt-app-operator:ci --go-test-flags "-v"
	./script/operator_courier_verify.sh
