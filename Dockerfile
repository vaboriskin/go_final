FROM golang:1.25 AS build
WORKDIR /src
    
COPY go.mod go.sum ./
RUN go mod download
    
COPY . .
    
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server .
    
FROM ubuntu:latest
WORKDIR /app
    
COPY --from=build /out/server /app/server
COPY web /app/web
    
  
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
    
EXPOSE 7540
    
CMD ["/app/server"]