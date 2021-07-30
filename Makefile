.PHONY: plugins core run clean
clean:
	rm plugins/*.so
	rm selfbot
core:
	go build -o selfbot
plugins:
	go build -o plugins/ping.so -buildmode=plugin ./plugins/ping
	go build -o plugins/ud.so -buildmode=plugin ./plugins/ud
	go build -o plugins/wiki.so -buildmode=plugin ./plugins/wiki
	go build -o plugins/vote.so -buildmode=plugin ./plugins/vote
	go build -o plugins/peertube.so -buildmode=plugin ./plugins/peertube
	go build -o plugins/screenshot.so -buildmode=plugin ./plugins/screenshot
run: core plugins
	./selfbot