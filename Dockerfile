FROM golang:1.21-bullseye AS builder

# 设置构建参数的默认值
ARG GIT_TAG="unknown"
ARG GIT_COMMIT_LOG="unknown"
ARG BUILD_TIME="unknown"
ARG BUILD_GO_VERSION="github@action golang:1.21-bullseye"


WORKDIR /go/src/app
COPY . .

# 打印构建参数
RUN echo "GIT_TAG=${GIT_TAG}"
RUN echo "GIT_COMMIT_LOG=${GIT_COMMIT_LOG}"
RUN echo "BUILD_TIME=${BUILD_TIME}"
RUN echo "BUILD_GO_VERSION=${BUILD_GO_VERSION}"


# 设置 LDFlags 变量
ENV LDFLAGS=" \
    -X 'main.GitTag=${GIT_TAG}' \
    -X 'main.GitCommitLog=${GIT_COMMIT_LOG}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.BuildGoVersion=${BUILD_GO_VERSION}' \
"

RUN wget -O gf "https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH)" && chmod +x gf && ./gf install -y && rm ./gf
RUN gf build -e "-ldflags \"${LDFLAGS}\" "

FROM loads/alpine:3.8

LABEL maintainer="Hamster <liaolaixin@gmail.com>"

###############################################################################
#                                INSTALLATION
###############################################################################

# 设置固定的项目路径
ENV WORKDIR /app/main
COPY --from=builder /go/src/app/temp/release/linux_amd64/service $WORKDIR/service
# 添加应用可执行文件，并设置执行权限
RUN chmod +x $WORKDIR/service
# 增加端口绑定
EXPOSE 10401

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ["./service"]




