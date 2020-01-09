build:
	GOOS=linux go build -a -tags netgo -ldflags '-w' -o app
	sudo docker build -t gcr.io/alien-emblem-264603/say .
	rm app

push:
	sudo gcloud docker -- push gcr.io/alien-emblem-264603/say
