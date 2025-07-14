DOCKER_CLI := `echo ${DOCKER_CLI:-docker}`

pg-start:
    @{{DOCKER_CLI}} run -d --name postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 6432:6432 postgres:latest

pg-clean:
    @{{DOCKER_CLI}} stop postgres
    @{{DOCKER_CLI}} rm postgres

rabbitmq-start:
    @{{DOCKER_CLI}} run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin rabbitmq:3.12-management
    @{{DOCKER_CLI}} cp deps/rabbitmq_delayed_message_exchange-3.12.0.ez rabbitmq:/plugins/.
    @{{DOCKER_CLI}} exec rabbitmq rabbitmq-plugins enable rabbitmq_delayed_message_exchange
    @{{DOCKER_CLI}} restart rabbitmq

rabbitmq-clean:
    @{{DOCKER_CLI}} stop rabbitmq
    @{{DOCKER_CLI}} rm rabbitmq

emqx-start:
    @{{DOCKER_CLI}} run -d --name emqx -p 1883:1883 -p 8081:8081 -p 8083:8083 -p 8084:8084 emqx:5.8.6

emqx-clean:
    @{{DOCKER_CLI}} stop emqx
    @{{DOCKER_CLI}} rm emqx
