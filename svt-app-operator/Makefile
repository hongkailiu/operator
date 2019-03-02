.PHONY : ci-test
ci-test:
	echo "ci-test for svt-app-operator"
	kubectl create namespace operator-test
	operator-sdk test local ./test/e2e --namespace operator-test
