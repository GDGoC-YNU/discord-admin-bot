FROM golang:1.23-alpine
WORKDIR /root
COPY . .
RUN go build -o /app/main .
ENV SECRET_LOCATION=/secrets/settings.yaml
CMD [ "/app/main" ]
