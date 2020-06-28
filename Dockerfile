FROM golang:1.14
WORKDIR /stars/
COPY ./ .
RUN  go build /stars/cmd/stars
CMD ["/stars/stars"]  