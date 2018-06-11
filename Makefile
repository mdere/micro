build_deploy:
	local_ensure.sh
	docker build -t mdere_micro:latest .
	docker tag mdere_micro:latest mdere/micro:latest
	docker push mdere/micro:latest

deploy:
	# docker tag mdere_micro:latest mdere/micro:latest
	

