.PHONY: pull-tags


VERSION:=`git describe --tags`
BUILD:=`date +%F`

pull-tags:
	git fetch --tags


gencode:
	echo Generating code
	for dir in $(shell find ../ -type d | grep -iv git); do \
		go generate $$dir ; \
	done


write-version:
	echo "package ham\n" > ../version.go
	echo "const Version = \"${VERSION}\"" >> ../version.go
	echo "const BuildDate = \"${BUILD}\"" >> ../version.go
