FROM ruby:3.1.2

RUN gem install fluentd --no-document

RUN gem install fluent-plugin-elasticsearch --no-document
RUN gem install fluent-plugin-kafka --no-document

RUN gem update fluentd
RUN gem update fluent-plugin-elasticsearch
RUN gem update fluent-plugin-kafka

COPY fluent.conf /fluentd/etc/fluent.conf

CMD ["fluentd", "-c", "/fluentd/etc/fluent.conf"]