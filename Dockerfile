FROM registry-vpc.cn-beijing.aliyuncs.com/bbt-base/centos:7.6.1810

COPY ./run /var/www/run

EXPOSE 9001