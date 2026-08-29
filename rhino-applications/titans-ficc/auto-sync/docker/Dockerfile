FROM docker2.gf.com.cn/library/ubuntu:18.04-latest

# 安装必要软件
RUN apt-get update && apt-get install -y \
    lsyncd \
    rsync \
    openssh-server \
    openssh-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 创建必要的运行目录和日志目录
RUN mkdir -p /var/log/lsyncd /run/sshd /var/run/sshd

# === 关键修复：在构建镜像时生成共享的SSH密钥对 ===
RUN mkdir -p /root/.ssh && \
    ssh-keygen -t rsa -N "" -f /root/.ssh/id_rsa && \
    # 将公钥复制为授权密钥，允许使用相同密钥的容器互相连接
    cp /root/.ssh/id_rsa.pub /root/.ssh/authorized_keys && \
    # 设置严格的权限，这是SSH安全所必需的
    chmod 600 /root/.ssh/id_rsa /root/.ssh/authorized_keys && \
    chmod 700 /root/.ssh

# 配置SSH服务器：允许root使用密钥登录，禁止密码登录
RUN sed -i 's/#PermitRootLogin.*/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config

# 复制启动脚本到镜像中
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# 声明容器暴露的端口
EXPOSE 22

# 设置启动脚本为入口点
ENTRYPOINT ["/entrypoint.sh"]