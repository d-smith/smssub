bin:
	GOOS=linux go build -o main
	zip smssub-deployment.zip main

clean:
	rm -f smssub-deployment.zip
