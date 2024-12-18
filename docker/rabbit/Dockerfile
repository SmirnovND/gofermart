FROM rabbitmq:3.12-management

# Устанавливаем необходимые утилиты и загружаем плагин
RUN apt-get update && apt-get install -y curl \
    && curl -L -o /plugins/rabbitmq_delayed_message_exchange.ez \
       https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.12.0/rabbitmq_delayed_message_exchange-3.12.0.ez \
    && rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange
