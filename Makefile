
run:
	# build new binary
	go build .

	# allow to use 80 port for non-root user
	sudo setcap 'cap_net_bind_service=+ep' touhou-music.online
	
	# run app
	./touhou-music.online

run-win:
	# build new binary
	go build .

	# run app
	./touhou-music.online

deploy:
	# building web service
	GOOS=linux GOARCH=amd64 go build .

	# archiving files required for deploy
	tar -czvf touhou-music.online.tar.gz --exclude static/music static \
		template touhou-music.online

	# copying deploy archive
	scp touhou-music.online.tar.gz root@95.213.237.2:/home/touhou-music.online/

	# extracting new files to working directory
	ssh root@95.213.237.2 tar -C /home/touhou-music.online/ -xf \
		/home/touhou-music.online/touhou-music.online.tar.gz

	# adding execute permissions
	ssh root@95.213.237.2 chmod +x \
		/home/touhou-music.online/touhou-music.online

	# allowing to bind port 80 for non-root user
	ssh root@95.213.237.2 setcap 'cap_net_bind_service=+ep' \
		/home/touhou-music.online/touhou-music.online

	# restarting web service for applying updates
	ssh root@95.213.237.2 systemctl restart touhou-music.online.service

logs:
	ssh root@95.213.237.2 tail -n100 -f \
		/var/log/touhou-music.online/touhou-music.online.log
