# Users Service

In order to test Github Actions and some AWS features, I created this simple users service

The application is Dockerized, has a docker-compose file to boot the service quickly and a Makefile to start working ASAP.

	make init
	make build
	make up

And you are ready to Go!

`init` will populate the `.env` file needed for injecting environment variables.  
`build` will create the development image to code inside of it.  
`up` will run the API, exposing ports specified in the docker-compose file.  

This lets the developer focus on the code, running it inside the container resembling production.

---

## Unit Testing

The unit testing was done using the dependency injection technique.

Same as the development, the unit testing is performed inside the docker container, to do so run the following:

	make devshell
	make t

`devshell` will run the development container and start a terminal inside of it.  
`t` will run the unit testing and provide the coverage level for each package.  



