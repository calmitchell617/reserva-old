FROM ubuntu:22.04

COPY /bin/linux_arm64/api /api

ARG DB_DSN=""
ENV DB_DSN="${DB_DSN}"

CMD sleep 3 \
  && /api -db-dsn=${DB_DSN} -smtp-host=${SMTP_HOST} -smtp-port=${SMTP_PORT} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD} -smtp-sender=${SMTP_SENDER}