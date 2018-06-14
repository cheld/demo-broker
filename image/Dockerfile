FROM ubuntu

RUN apt-get update &&  apt-get install -y ca-certificates

ADD servicebroker /opt/servicebroker/servicebroker
CMD /opt/servicebroker/servicebroker --help
