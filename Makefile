.PHONY: current all current-v8 current-arm64-executable android-executable-armv7 android-executable-arm64 android-executable-x86_64 android-executable-x86 windows-executable windows-executable-x86 windows-executable-x86_64 freebsd-executable freebsd-executable-x86 freebsd-executable-x86_64 freebsd-executable-arm64 netbsd-executable netbsd-executable-x86 netbsd-executable-x86_64 netbsd-executable-arm64 netbsd-executable netbsd-executable-x86 netbsd-executable-x86_64 netbsd-executable-arm64 openwrt-mt7620-mipsel_24kc
TARGETS:=build/ current current-no-readline current-v8
PACKAGETARGETS:=
ifneq ($(wildcard ${HOME}/android-ndk-r20b),)
	TARGETS:=${TARGETS} android-v8-executable-arm64 android-executable-armv7 android-executable-arm64 android-executable-x86_64 android-executable-x86
	PACKAGETARGETS:=${PACKAGETARGETS} package/android
endif
ifneq ($(wildcard /usr/bin/i686-w64-mingw32-gcc),)
	TARGETS:=${TARGETS} windows-executable-x86
endif
ifneq ($(wildcard /usr/bin/x86_64-w64-mingw32-gcc),)
	TARGETS:=${TARGETS} windows-executable-x86_64
endif
ifneq ($(wildcard /usr/bin/aarch64-linux-gnu-gcc),)
	TARGETS:=${TARGETS} current-arm64-executable
endif

VERSION=$(shell cat version)

SRCS_GO := $(foreach dir, $(shell find . -type d), $(wildcard $(dir)/*.go $(dir)/*.c))

CGO_DEF := "-DFB_VERSION=\"$(VERSION)\" -DFB_COMMIT=\"$(shell git log -1 --format=format:"%h")\" -DFB_COMMIT_LONG=\"$(shell git log -1 --format=format:"%H")\""

current: build/phoenixbuilder
all: ${TARGETS} build/hashes.json
current-no-readline: build/phoenixbuilder-no-readline
current-debug: build/phoenixbuilder-debug
current-v8: build/phoenixbuilder-v8
current-arm64-executable: build/phoenixbuilder-aarch64
android-executable-armv7: build/phoenixbuilder-android-static-executable-armv7 build/phoenixbuilder-android-termux-shared-executable-armv7 build/phoenixbuilder-android-shared-executable-armv7
android-executable-arm64: build/phoenixbuilder-android-static-executable-arm64 build/phoenixbuilder-android-termux-shared-executable-arm64 build/phoenixbuilder-android-shared-executable-arm64
android-v8-executable-arm64: build/phoenixbuilder-v8-android-static-executable-arm64 build/phoenixbuilder-v8-android-termux-shared-executable-arm64 build/phoenixbuilder-v8-android-shared-executable-arm64
android-executable-x86_64: build/phoenixbuilder-android-static-executable-x86_64 build/phoenixbuilder-android-termux-shared-executable-x86_64 build/phoenixbuilder-android-shared-executable-x86_64
android-executable-x86: build/phoenixbuilder-android-shared-executable-x86 build/phoenixbuilder-android-termux-shared-executable-x86 build/phoenixbuilder-android-static-executable-x86
windows-executable: windows-executable-x86 windows-executable-x86_64
windows-executable-x86: build/phoenixbuilder-windows-executable-x86.exe
windows-executable-x86_64: build/phoenixbuilder-windows-executable-x86_64.exe

package: ${PACKAGETARGETS}
release/:
	mkdir -p release
build/:
	mkdir build

#ifeq ($(shell uname | grep -iq 'Linux' && echo 1),1)
#ifeq ($(shell uname -m | grep -iqE "x86_64|amd64" && echo 1),1)

build/phoenixbuilder: build/ ${SRCS_GO}
	cd depends/stub&&make clean&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_ENABLED=1 go build -tags "${APPEND_GO_TAGS}" -trimpath -ldflags "-s -w" -o $@
build/phoenixbuilder-no-readline: build/ ${SRCS_GO}
	cd depends/stub&&make clean&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_ENABLED=1 go build -tags "no_readline" -tags "${APPEND_GO_TAGS}" -trimpath -ldflags "-s -w" -o $@
build/phoenixbuilder-with-symbols: build/ ${SRCS_GO}
	cd depends/stub&&make clean&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_ENABLED=1 go build -tags "${APPEND_GO_TAGS}" -trimpath -o $@
build/phoenixbuilder-v8: build/ ${SRCS_GO}
	cd depends/stub&&make clean&&cd -
	CGO_CFLAGS=${CGO_DEF}" -DWITH_V8" CGO_ENABLED=1 go build -tags "with_v8 ${APPEND_GO_TAGS}" -trimpath -ldflags "-s -w" -o build/phoenixbuilder-v8
build/libexternal_functions_provider.so: build/ io/external_functions_provider/provider.c
	gcc -shared io/external_functions_provider/provider.c -o build/libexternal_functions_provider.so
build/phoenixbuilder-static.a: build/ build/libexternal_functions_provider.so ${SRCS_GO}
	CGO_CFLAGS=${CGO_DEF} CGO_LDFLAGS="-Lbuild -lexternal_functions_provider" CGO_ENABLED=1 go build -trimpath -buildmode=c-archive -ldflags "-s -w" -tags "no_readline,is_tweak ${APPEND_GO_TAGS}" -o build/phoenixbuilder-static.a
build/phoenixbuilder-aarch64: build/ ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=/usr/bin/aarch64-linux-gnu-gcc&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=/usr/bin/aarch64-linux-gnu-gcc CGO_ENABLED=1 GOARCH=arm64 go build -tags use_aarch64_linux_rl -trimpath -ldflags "-s -w" -o build/phoenixbuilder-aarch64
build/phoenixbuilder-android-static-executable-armv7: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang GOOS=android GOARCH=arm CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-static-executable-armv7
build/phoenixbuilder-v8-android-static-executable-arm64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF}" -DWITH_V8" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -tags with_v8 -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o build/phoenixbuilder-v8-android-static-executable-arm64
build/phoenixbuilder-android-static-executable-arm64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o build/phoenixbuilder-android-static-executable-arm64
build/phoenixbuilder-android-static-executable-x86: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang GOOS=android GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-static-executable-x86
build/phoenixbuilder-android-static-executable-x86_64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang GOOS=android GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-static-executable-x86_64
build/phoenixbuilder-android-termux-shared-executable-armv7: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang GOOS=android GOARCH=arm CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-termux-shared-executable-armv7
build/phoenixbuilder-v8-android-termux-shared-executable-arm64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF}" -DWITH_V8" CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -tags with_v8,android_shared -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o build/phoenixbuilder-v8-android-termux-shared-executable-arm64
build/phoenixbuilder-android-termux-shared-executable-arm64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o build/phoenixbuilder-android-termux-shared-executable-arm64
build/phoenixbuilder-android-termux-shared-executable-x86: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang GOOS=android GOARCH=386 CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-termux-shared-executable-x86
build/phoenixbuilder-android-termux-shared-executable-x86_64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang GOOS=android GOARCH=amd64 CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-termux-shared-executable-x86_64
build/phoenixbuilder-android-executable-arm64: build/phoenixbuilder-android-termux-shared-executable-arm64
	ln build/phoenixbuilder-android-termux-shared-executable-arm64 build/phoenixbuilder-android-executable-arm64
build/phoenixbuilder-android-executable-armv7: build/phoenixbuilder-android-termux-shared-executable-armv7
	ln build/phoenixbuilder-android-termux-shared-executable-armv7 build/phoenixbuilder-android-executable-armv7
build/phoenixbuilder-android-executable-x86: build/phoenixbuilder-android-termux-shared-executable-x86
	ln build/phoenixbuilder-android-termux-shared-executable-x86 build/phoenixbuilder-android-executable-x86
build/phoenixbuilder-android-executable-x86_64: build/phoenixbuilder-android-termux-shared-executable-x86_64
	ln build/phoenixbuilder-android-termux-shared-executable-x86_64 build/phoenixbuilder-android-executable-x86_64
build/phoenixbuilder-android-shared-executable-armv7: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi16-clang GOOS=android GOARCH=arm CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-shared-executable-armv7
build/phoenixbuilder-v8-android-shared-executable-arm64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF}" -DWITH_V8" CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -tags with_v8,android_shared -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o build/phoenixbuilder-v8-android-shared-executable-arm64
build/phoenixbuilder-android-shared-executable-arm64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang CXX=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o build/phoenixbuilder-android-shared-executable-arm64
build/phoenixbuilder-android-shared-executable-x86: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android21-clang GOOS=android GOARCH=386 CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-shared-executable-x86
build/phoenixbuilder-android-shared-executable-x86_64: build/ ${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang ${SRCS_GO}
	cd depends/stub&&make clean&&make CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang&&cd -
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=${HOME}/android-ndk-r20b/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android21-clang GOOS=android GOARCH=amd64 CGO_ENABLED=1 go build -tags android_shared -trimpath -ldflags "-s -w" -o build/phoenixbuilder-android-shared-executable-x86_64
build/phoenixbuilder-windows-executable-x86.exe: build/ /usr/bin/i686-w64-mingw32-gcc ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=/usr/bin/i686-w64-mingw32-gcc GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o build/phoenixbuilder-windows-executable-x86.exe
build/phoenixbuilder-windows-executable-x86_64.exe: build/ /usr/bin/x86_64-w64-mingw32-gcc ${SRCS_GO}
	GODEBUG=madvdontneed=1 CGO_CFLAGS=${CGO_DEF} CC=/usr/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o build/phoenixbuilder-windows-executable-x86_64.exe
build/hashes.json: build genhash.js ${TARGETS}
	node genhash.js
	cp version build/version

package/android: package/android-armv7 package/android-arm64
package/android-armv7: build/phoenixbuilder-android-executable-armv7 release/
	mkdir -p release/phoenixbuilder-android-armv7/data/data/com.termux/files/usr/bin release/phoenixbuilder-android-armv7/DEBIAN
	cp build/phoenixbuilder-android-executable-armv7 release/phoenixbuilder-android-armv7/data/data/com.termux/files/usr/bin/fastbuilder
	printf "Package: pro.fastbuilder.phoenix-android\n\
	Name: FastBuilder Phoenix (Alpha)\n\
	Version: $(VERSION)\n\
	Architecture: arm\n\
	Depends: libreadline8 | readline (>= 8.0.0), zlib | zlib1g\n\
	Maintainer: Ruphane\n\
	Author: Bouldev <admin@boul.dev>\n\
	Section: Games\n\
	Priority: optional\n\
	Homepage: https://fastbuilder.pro\n\
	Description: Modern Minecraft structuring tool\n" > release/phoenixbuilder-android-armv7/DEBIAN/control
	dpkg-deb -Zxz -b release/phoenixbuilder-android-armv7 release/
package/android-arm64: build/phoenixbuilder-android-executable-arm64 release/
	mkdir -p release/phoenixbuilder-android-arm64/data/data/com.termux/files/usr/bin release/phoenixbuilder-android-arm64/DEBIAN
	cp build/phoenixbuilder-android-executable-arm64 release/phoenixbuilder-android-arm64/data/data/com.termux/files/usr/bin/fastbuilder
	printf "Package: pro.fastbuilder.phoenix-android\n\
	Name: FastBuilder Phoenix (Alpha)\n\
	Version: $(VERSION)\n\
	Architecture: aarch64\n\
	Depends: libreadline8 | readline (>= 8.0.0), zlib | zlib1g\n\
	Maintainer: Ruphane\n\
	Author: Bouldev <admin@boul.dev>\n\
	Section: Games\n\
	Priority: optional\n\
	Homepage: https://fastbuilder.pro\n\
	Description: Modern Minecraft structuring tool\n" > release/phoenixbuilder-android-arm64/DEBIAN/control
	dpkg-deb -Zxz -b release/phoenixbuilder-android-arm64 release/
clean:
	rm -f build/phoenixbuilder*
