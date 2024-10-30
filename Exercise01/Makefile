all:
	unzip Additional.zip
	rm -rf __MACOSX
	mkdir myblog-binary
	cd pkg && go build
	mv pkg/myblog myblog-binary/
	cp pkg/admin_credentials.txt myblog-binary/
	cd myblog-binary && ./myblog

clean:
	rm -rf __MACOSX images templates myblog-binary

rest: clean all