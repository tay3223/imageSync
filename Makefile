######################################
# 全局变量
######################################
#要编译的命令名称
NAME := imageSync
#版本
VERSION := 0.0.6
#编译输出目录
OUTPUT_PATH := ./build/${VERSION}/
#是否开启cgo（0代表不开启，1代表开始）
CGO_STATUS := 0


######################################
# 局部变量
######################################
build_mac: baoName := ${NAME}-${VERSION}-mac-amd64
build_linux: baoName := ${NAME}-${VERSION}-linux-amd64
build_windows: baoName := ${NAME}-${VERSION}-windows-amd64




#指定加载哪些伪造的Target
.PHONY: clean  start  build-linux  build_windows  build_mac
#指定缺省状态下执行哪些Target
all: clean  start  build_linux  build_windows  build_mac


#Target：主要用来清理一些开发和编译过程中的无用文件
clean:
	rm -rf ./build/*
	rm -rf ./*.log
	rm -rf ./*.tar.gz
	rm -rf ./.imageSync


#Target：项目依赖
start:
	go mod tidy


#Target：交叉编译到Linux平台
build_linux:
	mkdir -p ${OUTPUT_PATH}${baoName}/
	GOARCH=amd64 \
	GOOS=linux \
	CGO_ENABLED=${CGO_STATUS} \
	CGO_LDFLAGS="-static" \
	CC=x86_64-linux-musl-gcc \
	CXX=x86_64-linux-musl-g++ \
	go build -o "${OUTPUT_PATH}${baoName}/${NAME}"
	tar -zcvf ${OUTPUT_PATH}${baoName}.tar.gz -C ${OUTPUT_PATH} ${baoName}


#Target：交叉编译到windows平台
build_windows:
	mkdir -p ${OUTPUT_PATH}${baoName}/
	GOARCH=amd64 \
	GOOS=windows \
	CGO_ENABLED=${CGO_STATUS} \
	CGO_CFLAGS="-g -O2 -Wno-return-local-addr" \
	CC=x86_64-w64-mingw32-gcc \
	CXX=x86_64-w64-mingw32-g++ \
	go build -o "${OUTPUT_PATH}${baoName}/${NAME}"
	tar -zcvf ${OUTPUT_PATH}${baoName}.tar.gz -C ${OUTPUT_PATH} ${baoName}


#Target：交叉编译到Mac平台
build_mac:
	mkdir -p ${OUTPUT_PATH}${baoName}/
	GOARCH=amd64 \
	GOOS=darwin \
	CGO_ENABLED=${CGO_STATUS} \
	go build -o "${OUTPUT_PATH}${baoName}/${NAME}"
	tar -zcvf ${OUTPUT_PATH}${baoName}.tar.gz -C ${OUTPUT_PATH} ${baoName}
	cp -f "${OUTPUT_PATH}${baoName}/${NAME}" /usr/local/bin/
	chmod +x "${OUTPUT_PATH}${baoName}/${NAME}"
	chmod +x "/usr/local/bin/${NAME}"
