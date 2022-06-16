FROM rabbitmq:3.10.5-alpine
EXPOSE 5672
VOLUME /var/lib/rabbitmq
COPY rabbitmq.conf /etc/rabbitmq/rabbitmq.conf
RUN chown rabbitmq:rabbitmq /etc/rabbitmq/rabbitmq.conf