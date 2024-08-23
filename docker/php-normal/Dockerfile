FROM php:7.4.33-apache
ENV TZ=Asia/Shanghai
RUN echo '' > /etc/apt/sources.list
RUN sed -i "1ideb http://mirrors.aliyun.com/debian/ bullseye main non-free contrib" /etc/apt/sources.list
RUN sed -i "2ideb-src http://mirrors.aliyun.com/debian/ bullseye main non-free contrib" /etc/apt/sources.list
RUN sed -i "3ideb http://mirrors.aliyun.com/debian-security/ bullseye-security main" /etc/apt/sources.list
RUN sed -i "4ideb-src http://mirrors.aliyun.com/debian-security/ bullseye-security main" /etc/apt/sources.list
RUN sed -i "5ideb http://mirrors.aliyun.com/debian/ bullseye-updates main non-free contrib" /etc/apt/sources.list
RUN sed -i "6ideb-src http://mirrors.aliyun.com/debian/ bullseye-updates main non-free contrib" /etc/apt/sources.list
RUN sed -i "7ideb http://mirrors.aliyun.com/debian/ bullseye-backports main non-free contrib" /etc/apt/sources.list
RUN sed -i "8ideb-src http://mirrors.aliyun.com/debian/ bullseye-backports main non-free contrib" /etc/apt/sources.list

RUN a2enmod rewrite

ADD https://github.com/mlocati/docker-php-extension-installer/releases/latest/download/install-php-extensions /usr/local/bin/
RUN chmod +x /usr/local/bin/install-php-extensions

RUN install-php-extensions pdo_mysql pdo_pgsql xdebug intl pcntl gd zip imagick memcache @composer
