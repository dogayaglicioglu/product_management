FROM fluent/fluentd:v1.16.2-debian-1.0

USER root

RUN gem install elasticsearch -v '7.17.0'
RUN gem install fluent-plugin-elasticsearch -v '5.0.3'
RUN apt-get update && apt-get install -y curl

USER fluent

COPY fluent.conf /fluentd/etc/