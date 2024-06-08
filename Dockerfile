FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/dops .


FROM alpine AS runner

RUN apk add --no-cache -u tzdata

ENV TZ=Asia/Seoul

WORKDIR /app

COPY --from=builder /app/dops .

CMD ["./dops"]
