build:
	docker build . -t dnsp
serve:
	docker run --rm -p 53:53/udp dnsp --resolve 8.8.4.4,8.8.8.8 --blacklist=https://hosts-file.net/download/hosts.txt --poll 1h -l 0.0.0.0:53
debug:
	docker run --rm -p 5353:53/udp dnsp --resolve 8.8.4.4,8.8.8.8 --blacklist=https://hosts-file.net/download/hosts.txt --poll 1h -l 0.0.0.0:53
