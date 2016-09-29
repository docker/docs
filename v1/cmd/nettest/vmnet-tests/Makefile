.PHONY: all

all:
	xcodebuild -project nettest.xcodeproj

run: all
	sudo ./build/Release/nettest

clean: 
	rm -rf ./build
