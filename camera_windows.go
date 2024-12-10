package camera

/*
#cgo LDFLAGS: -lole32 -lmf -lmfplat -lmfuuid -lmfreadwrite

#define COBJMACROS
#include <windows.h>
#include <mfapi.h>
#include <mfidl.h>
#include <mfreadwrite.h>
#include <mfobjects.h>
#include <stdint.h>
#include <stdio.h>

#define CHECK_HR(hr) if (FAILED(hr)) { fprintf(stderr, "Error: 0x%08X, at line %d\n", hr, __LINE__); return -1; }

#define HI32(x) ((UINT32)(((x) >> 32) & 0xFFFFFFFF))
#define LO32(x) ((UINT32)((x) & 0xFFFFFFFF))

// Declare the callback function
extern void onFrameAvailableGo(void* data, int width, int height, int bytesPerPixel);

static IMFSourceReader *pReader_ = NULL;
static UINT32 width_ = 0;
static UINT32 height_ = 0;
static GUID subtype_ = {0};  // Declare the subtype variable
static uint8_t *rgbData_ = NULL;

// Function to convert YUY2 to RGB24
static void YUY2toRGB24(uint8_t *yuy2, uint8_t *rgb, int width, int height) {
    for (int i = 0; i < width * height * 2; i += 4) {
        uint8_t y0 = yuy2[i];
        uint8_t u = yuy2[i + 1];
        uint8_t y1 = yuy2[i + 2];
        uint8_t v = yuy2[i + 3];

        int c = y0 - 16;
        int d = u - 128;
        int e = v - 128;

        rgb[(i / 2) * 3] = (uint8_t) (1.164 * c + 1.596 * e);
        rgb[(i / 2) * 3 + 1] = (uint8_t) (1.164 * c - 0.392 * d - 0.813 * e);
        rgb[(i / 2) * 3 + 2] = (uint8_t) (1.164 * c + 2.017 * d);

        c = y1 - 16;

        rgb[(i / 2 + 1) * 3] = (uint8_t) (1.164 * c + 1.596 * e);
        rgb[(i / 2 + 1) * 3 + 1] = (uint8_t) (1.164 * c - 0.392 * d - 0.813 * e);
        rgb[(i / 2 + 1) * 3 + 2] = (uint8_t) (1.164 * c + 2.017 * d);
    }
}

static int webcam_open(int deviceId, int width, int height) {
    HRESULT hr;
    IMFActivate **ppDevices = NULL;
    IMFAttributes *pAttributes = NULL;
    UINT32 count = 0;

    width_ = width;
    height_ = height;

    hr = MFStartup(MF_VERSION, MFSTARTUP_NOSOCKET);
    CHECK_HR(hr);

    hr = MFCreateAttributes(&pAttributes, 1);
    CHECK_HR(hr);

    hr = IMFAttributes_SetGUID(pAttributes, &MF_DEVSOURCE_ATTRIBUTE_SOURCE_TYPE, &MF_DEVSOURCE_ATTRIBUTE_SOURCE_TYPE_VIDCAP_GUID);
    CHECK_HR(hr);

    hr = MFEnumDeviceSources(pAttributes, &ppDevices, &count);
    CHECK_HR(hr);

    if (deviceId >= (int)count) {
        fprintf(stderr, "Error: Device ID %d out of range, %u devices found\n", deviceId, count);
        IMFAttributes_Release(pAttributes);
        CoTaskMemFree(ppDevices);
        return -1;
    }

    IMFMediaSource *pSource = NULL;
    hr = IMFActivate_ActivateObject(ppDevices[deviceId], &IID_IMFMediaSource, (void**)&pSource);
    CHECK_HR(hr);

    hr = MFCreateSourceReaderFromMediaSource(pSource, NULL, &pReader_);
    CHECK_HR(hr);

    IMFMediaSource_Release(pSource);

    IMFAttributes_Release(pAttributes);
    for (UINT32 i = 0; i < count; i++) {
        IMFActivate_Release(ppDevices[i]);
    }
    CoTaskMemFree(ppDevices);

	rgbData_ = (uint8_t*) malloc(width * height * 4);
    //printf("Webcam opened successfully.\n");
    return 0;
}

// Function to print GUID
static void PrintGUID(const GUID *guid) {
    WCHAR guidString[39] = {0};
    StringFromGUID2(guid, guidString, sizeof(guidString) / sizeof(WCHAR));
    wprintf(L"%ls", guidString);
}

static int webcam_start() {
    HRESULT hr;
    IMFMediaType *pNativeType = NULL;
    UINT32 i = 0;
    UINT32 bestWidth = 0, bestHeight = 0;

    while (SUCCEEDED(hr = IMFSourceReader_GetNativeMediaType(pReader_, MF_SOURCE_READER_FIRST_VIDEO_STREAM, i, &pNativeType))) {
        GUID majorType = {0};
        hr = IMFMediaType_GetGUID(pNativeType, &MF_MT_MAJOR_TYPE, &majorType);
        if (FAILED(hr)) {
            IMFMediaType_Release(pNativeType);
            CHECK_HR(hr);
        }

        if (IsEqualGUID(&majorType, &MFMediaType_Video)) {
            UINT64 frameSize = 0;
            hr = IMFMediaType_GetUINT64(pNativeType, &MF_MT_FRAME_SIZE, &frameSize);
            CHECK_HR(hr);

            UINT32 width = HI32(frameSize);
            UINT32 height = LO32(frameSize);

            if (width == width_ && height == height_) {
                hr = IMFMediaType_GetGUID(pNativeType, &MF_MT_SUBTYPE, &subtype_);
                CHECK_HR(hr);

                hr = IMFSourceReader_SetCurrentMediaType(pReader_, MF_SOURCE_READER_FIRST_VIDEO_STREAM, NULL, pNativeType);
                IMFMediaType_Release(pNativeType);
                CHECK_HR(hr);
                //printf("Selected resolution: %ux%u\n", width, height);
                //printf("Selected video subtype: ");
                //PrintGUID(&subtype_);
                return 0;
            }

            if ((bestWidth == 0 && bestHeight == 0) || (width <= width_ && height <= height_)) {
                bestWidth = width;
                bestHeight = height;
                hr = IMFMediaType_GetGUID(pNativeType, &MF_MT_SUBTYPE, &subtype_);
                CHECK_HR(hr);
            }
        }

        IMFMediaType_Release(pNativeType);
        i++;
    }

    if (bestWidth == 0 || bestHeight == 0) {
        fprintf(stderr, "Failed to find a suitable media type.\n");
        return -1;
    }

    width_ = bestWidth;
    height_ = bestHeight;
    hr = IMFSourceReader_SetCurrentMediaType(pReader_, MF_SOURCE_READER_FIRST_VIDEO_STREAM, NULL, pNativeType);
    CHECK_HR(hr);

    //printf("Selected resolution: %ux%u\n", width_, height_);
    //printf("Selected video subtype: ");
    //PrintGUID(&subtype_);
    //printf("\n");

    return 0;
}

static int capture_frame() {
    HRESULT hr;
    IMFSample *pSample = NULL;
    DWORD streamIndex, flags;
    LONGLONG timestamp;

    hr = IMFSourceReader_ReadSample(pReader_, MF_SOURCE_READER_FIRST_VIDEO_STREAM, 0, &streamIndex, &flags, &timestamp, &pSample);
    CHECK_HR(hr);

    if (!pSample) {
		return -1;
    }

	IMFMediaBuffer *pBuffer = NULL;
	hr = IMFSample_ConvertToContiguousBuffer(pSample, &pBuffer);
	CHECK_HR(hr);

	BYTE *pData = NULL;
	DWORD maxLength = 0, currentLength = 0;
	hr = IMFMediaBuffer_Lock(pBuffer, &pData, &maxLength, &currentLength);
	CHECK_HR(hr);

	if (IsEqualGUID(&subtype_, &MFVideoFormat_YUY2)) {
		// Handle YUY2 format appropriately
		if (rgbData_ == NULL) {
			fprintf(stderr, "Failed to allocate memory for RGB data\n");
			IMFMediaBuffer_Unlock(pBuffer);
			IMFMediaBuffer_Release(pBuffer);
			IMFSample_Release(pSample);
			return -1;
		}

		YUY2toRGB24(pData, rgbData_, width_, height_);
		onFrameAvailableGo(rgbData_, width_, height_, 3);
	} else if (IsEqualGUID(&subtype_, &MFVideoFormat_RGB24)) {
		// Handling RGB24 format
		if (maxLength < width_ * height_ * 3 || currentLength < width_ * height_ * 3) {
			fprintf(stderr, "Buffer size is too small for the expected frame size\n");
		} else {
			onFrameAvailableGo(pData, width_, height_, 3);
		}
	} else {
		// Handling other formats
	}

	hr = IMFMediaBuffer_Unlock(pBuffer);
	IMFMediaBuffer_Release(pBuffer);
	IMFSample_Release(pSample);
	CHECK_HR(hr);

    return 0;
}

static int webcam_stop() {
    if (pReader_) {
        IMFSourceReader_Release(pReader_);
        pReader_ = NULL;
    }
    MFShutdown();
    return 0;
}

static void webcam_delete() {
    if (pReader_) {
        IMFSourceReader_Release(pReader_);
        pReader_ = NULL;
    }
	MFShutdown();
	free(rgbData_);
	rgbData_ = NULL;
}

static void copyImage(uint8_t *dstBuf, void* srcBuf, size_t frame_size) {
    memcpy(dstBuf, srcBuf, frame_size);
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

var (
	running = false
)

func startCapture() {
	if running {
		return
	}

	running = true
	go func() {
		for running {
			C.capture_frame()
		}
		running = false
	}()
}

//export onFrameAvailableGo
func onFrameAvailableGo(data unsafe.Pointer, width, height, bytesPerPixel C.int) {
	frameSize := int(width) * int(height) * int(bytesPerPixel)
	go func() {
		buf := make([]byte, frameSize)
		C.copyImage((*C.uint8_t)(unsafe.Pointer(&buf[0])), data, C.size_t(frameSize))

		// Convert the buffer to an image.RGBA
		rgba := convertRGB24ToRGBA(buf, int(width), int(height))

		select {
		case frameBufferChan <- rgba:
		default:
			// Drop the frame if the channel is full
		}
	}()
}

func openCamera(id, width, height int) error {
	if C.webcam_open(C.int(id), C.int(width), C.int(height)) != 0 {
		return fmt.Errorf("failed to initialize camera")
	}
	return nil
}

func startCamera() error {
	if C.webcam_start() != 0 {
		return fmt.Errorf("failed to start camera")
	}

	startCapture()
	return nil
}

func stopCamera() error {
	running = false
	if C.webcam_stop() != 0 {
		return fmt.Errorf("failed to stop camera")
	}
	return nil
}

func closeCamera() {
	C.webcam_delete()
}
