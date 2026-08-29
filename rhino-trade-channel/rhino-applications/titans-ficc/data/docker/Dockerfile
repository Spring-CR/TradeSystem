FROM docker2.gf.com.cn/library/ubuntu:18.04-latest
#FROM docker2.gf.com.cn/library/java:v1.0.0-oraclejdk8-maven-3.3.9
ADD ./data-titans-ficc /

# 配置目录并安装日志轮转工具
RUN mkdir -p /opt/tradeclient; mkdir -p /opt/working_dir; mkdir -p /opt/log; apt-get update && \
    apt-get install -y logrotate cron && \
    rm -rf /var/lib/apt/lists/*

# 配置日志轮转规则（每日00:00触发）
RUN echo "/opt/log/data.log {"  > /etc/logrotate.d/data && \
    echo "  daily"                     >> /etc/logrotate.d/data && \
    echo "  rotate 7"                  >> /etc/logrotate.d/data && \
    echo "  compress"                  >> /etc/logrotate.d/data && \
    echo "  delaycompress"             >> /etc/logrotate.d/data && \
    echo "  missingok"                 >> /etc/logrotate.d/data && \
    echo "  notifempty"                >> /etc/logrotate.d/data && \
    echo "  dateext"                   >> /etc/logrotate.d/data && \
    echo "  copytruncate"              >> /etc/logrotate.d/data && \
    echo "  create 644 root root"      >> /etc/logrotate.d/data && \
    echo "}"

# 设置定时任务（每日00:00执行）
RUN echo "0 0 * * * root /usr/sbin/logrotate -f /etc/logrotate.d/data" > /etc/cron.d/trade-logrotate && \
    chmod 644 /etc/cron.d/trade-logrotate


WORKDIR /opt/log
VOLUME /opt/log
CMD ["/bin/sh", "-c", "service cron start;/data-titans-ficc 2>&1 | tee -a /opt/log/data.log"]
