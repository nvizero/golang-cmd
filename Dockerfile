# Start by building the application.
#docker pull golang:1.18.10-bullseye
FROM golang:1.21.3-bullseye

WORKDIR /usr/local/go/src/udate
COPY . .
COPY ./known_hosts /root/.ssh/known_hosts
COPY ./id_rsa /root/.ssh/id_rsa
COPY ./id_rsa.pub /root/.ssh/id_rsa.pub
RUN chmod -R 600 /root/.ssh/id_rsa  
RUN chmod -R 600 /root/.ssh/id_rsa.pub  
RUN go mod init udate 
RUN go mod tidy 
RUN cd /usr/local/go/src/udate/control && rm go.*  
RUN cd /usr/local/go/src/udate/utils && rm go.*  
RUN go mod download
RUN CGO_ENABLED=0 go build -o /usr/local/go/src/udate

# Now copy it into our base image.
# FROM alpine:latest
# RUN apk add openssh
# WORKDIR /app
# COPY . .
# COPY --from=build /usr/local/go/src/udate/udate /app/udate
EXPOSE 8080
CMD ["/usr/local/go/src/udate/udate"]
