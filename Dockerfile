FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -tags 'awskms awssecretsmanager awsssm' -ldflags='-s -w' -o /exec-with-secrets .

FROM amazonlinux:2023
COPY --from=build /exec-with-secrets /exec-with-secrets

ENTRYPOINT ["/exec-with-secrets"]

CMD ["env"]
