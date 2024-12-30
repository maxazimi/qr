#
# Cross-platform QR scanner with Gio UI (macOS, Windows, Linux & Android are currently supported)
#
APP_NAME=qr
MODULE=github.com/maxazimi/$(APP_NAME)
APPID=com.maxazimi.$(APP_NAME)
DEBUG_APK=$(APP_NAME)-debug.apk
RELEASE_APK=$(APP_NAME)-release.apk
LIB_DIR=android/libs
AAR=$(LIB_DIR)/$(APP_NAME).aar
KEYSTORE=$(HOME)/.android/$(APP_NAME)-release.keystore

#
# Android specific
#
ARCH=arm64
MIN_SDK=24
SDK_ROOT=$(HOME)/Android/Android.SDK.Release.26.1.1.Mac
NDK_ROOT=$(SDK_ROOT)/ndk-bundle
ZIPALIGN=$(SDK_ROOT)/build-tools/28.0.2/zipalign
APKSIGNER=$(SDK_ROOT)/build-tools/28.0.2/apksigner

# Detect the operating system
OS := $(shell uname -s)

# Define the targets for each OS
ifeq ($(OS), Linux)
	OS_TARGET = linux
else ifeq ($(OS), Darwin)
	OS_TARGET = macos
else ifneq (,$(findstring NT,$(OS)))
	OS_TARGET = windows
else
	$(error Unsupported operating system: $(OS))
endif

all: $(OS_TARGET)

apk_install: $(DEBUG_APK)
	@adb install -r $<
	@adb shell am start -n $(APPID)/org.gioui.GioActivity

android: $(DEBUG_APK)

$(DEBUG_APK): $(AAR)
	@(cd android && ./gradlew test assembleDebug)
	@mv android/build/outputs/apk/debug/*-debug.apk $@

$(RELEASE_APK): $(AAR)
	@(cd android && ./gradlew assembleRelease)
	@$(ZIPALIGN) -v -p 4 android/build/outputs/apk/release/*-release-unsigned.apk $(APP_NAME)-unsigned-aligned.apk
	@$(APKSIGNER) sign --ks $(KEYSTORE) --out $@ $(APP_NAME)-unsigned-aligned.apk
	@rm $(APP_NAME)-unsigned-aligned.apk

aar: $(AAR)

$(AAR):
	@mkdir -p $(LIB_DIR)
	@go run gioui.org/cmd/gogio -buildmode archive -target android -arch $(ARCH) -minsdk $(MIN_SDK) \
		-appid $(APPID) -tags novulkan -o $(AAR) $(MODULE)

apk_genkey:
	keytool -genkey -v -keystore $(KEYSTORE) -keyalg RSA -keysize 2048 -validity 10000

macos:
	@(CGO_CFLAGS="-Wno-format" go build -ldflags '-s -w' -o $(APP_NAME))

linux:
	@(CGO_CFLAGS="-Wno-format" go build -ldflags '-s -w' -o $(APP_NAME))

windows:
	@(CGO_CFLAGS="-Wno-format" go build -ldflags '-s -w' -o $(APP_NAME).exe)

clean:
	@rm -f $(APP_NAME) $(RELEASE_APK) $(DEBUG_APK) *.exe
	@rm -rf $(APP_NAME).* $(APP_NAME)-* android/build $(LIB_DIR)/*

.PHONY: all clean install android macos linux windows aar $(AAR) $(DEBUG_APK) $(RELEASE_AAB)
