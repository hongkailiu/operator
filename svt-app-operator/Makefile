.PHONY : ci-test
ci-test:
	echo "ci-test for svt-app-operator"
	kubectl create namespace operator-test
	operator-sdk test local ./test/e2e --namespace operator-test
	operator-sdk build docker.io/hongkailiu/svt-app-operator:ci
	docker images
	./script/operator_courier_verify.sh
