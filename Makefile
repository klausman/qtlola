.PHONY: all clean win linux lopt wopt

all: linux

linux: qtlola

qtlola:
	go build -v -ldflags="-s -w"

lopt: qtlola
	upx qtlola

win: deploy/windows/qtlola.exe

deploy/windows/qtlola.exe:
	qtdeploy -docker -ldflags="-s -w" build windows_64_static

wopt: win
	upx deploy/windows/qtlola.exe

clean:
	rm -rf windows linux deploy qtlola
