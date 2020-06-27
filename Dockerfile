FROM golang

ADD /server .

RUN ls

RUN go install .

ENTRYPOINT /go/bin/realtimedashboard

EXPOSE 8080