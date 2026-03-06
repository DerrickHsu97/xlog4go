PKG := git.xiaojukeji.com/shield-arch/dlog4go
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ |grep -v example)

atest:
	go test -short ${PKG_LIST}

coverage:
	go test -covermode=count -v -coverprofile cover.cov ${PKG_LIST}

html:coverage
	go tool cover -html=cover.cov

clean:
	rm -f ./test/demo*
	rm -f ./test/public.log
	rm -f cover.cov
