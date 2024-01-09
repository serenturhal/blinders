golint:
	go list -f '{{.Dir}}/...' -m | xargs -n 1 sh -c 'golangci-lint run $$0 || exit -1'
