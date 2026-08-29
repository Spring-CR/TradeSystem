FROM docker2.gf.com.cn/library/python:3.8-latest
ADD sources.list /etc/apt/sources.list
ADD ./rhino-api-server /
WORKDIR /
CMD ["/rhino-api-server"]