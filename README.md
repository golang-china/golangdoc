# Godoc for Golang, support translate.

Install golangdoc:

	go get github.com/golang-china/golangdoc

Download Chinese Translate:

	git clone https://github.com/golang-china/golangdoc.translations.git $(GOROOT)/translations

Start Chinese Godoc Server:

	golangdoc -http=:6060 -lang=zh_CN

See:

- http://github.com/golang-china/golangdoc.translations
- http://godoc.org/github.com/golang-china/golangdoc/docgen
- http://godoc.org/github.com/golang-china/golangdoc/local


# Run as windows service:

	# install as windows service
	golangdoc -service-install -http=:6060

	# start/stop service
	golangdoc -service-start
	golangdoc -service-stop

	# remove service
	golangdoc -service-remove

# BUGS

Report bugs to chaishushan@gmail.com.

Thanks!
