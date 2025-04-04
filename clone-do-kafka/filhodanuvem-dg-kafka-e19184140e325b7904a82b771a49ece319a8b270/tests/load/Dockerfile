FROM golang:1.18.2

WORKDIR /kafka
# install xk6 (k6 plus extension to execute binaries)
RUN go install go.k6.io/xk6/cmd/xk6@v0.9.0
RUN xk6 build --with github.com/grafana/xk6-exec@v0.1.0
RUN chmod +x /kafka/k6

# compile kafka-client
COPY . .
RUN go build -o kafka-cli ./cmd/cli/main.go 
RUN mv kafka-cli /usr/local/bin/kafka-cli

ENTRYPOINT ["/kafka/k6"] 