FROM scratch

MAINTAINER Heifeng <0987363@gmail.com>

ADD vsub /

EXPOSE 10090
CMD ["/vsub", "serve"]
