default:
	go build
	go install
	sudo supervisorctl restart blog
