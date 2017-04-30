default:
	go install
	sudo supervisorctl restart blog
