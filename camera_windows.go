package camera

/*
#cgo LDFLAGS: -lvfw32
#include <windows.h>
#include <vfw.h>
#include <stdint.h>
#include <stdio.h>

#define FOURCC(a, b, c, d) ((DWORD)(((d) << 24) | ((c) << 16) | ((b) << 8) | (a)))

typedef struct {
    HWND hwnd;
    long w, h;
    BITMAPINFO bmi;
    void *rgb;
} CAM;

// Global device pointer
static unsigned long _frame_size = 0;
static CAM* _device = NULL;

static void copyImage(uint8_t *dstBuf, void* srcBuf) {
    memcpy(dstBuf, srcBuf, _frame_size);
}

static unsigned char clip(int i) {
    if (i <= 0)
        return 0;
    if (i >= 255)
        return 255;
    return (unsigned char)i;
}

static void YUV444toRGB888(unsigned char y, unsigned char u, unsigned char v, unsigned char *dst) {
    int C = y - 16;
    int D = u - 128;
    int E = v - 128;
    *dst++ = clip((298 * C + 409 * E + 128) >> 8);
    *dst++ = clip((298 * C - 100 * D - 208 * E + 128) >> 8);
    *dst++ = clip((298 * C + 516 * D + 128) >> 8);
}

static void YUYVToRGB24(int w, int h, unsigned char *src, unsigned char *dst) {
    int i;
    unsigned char u, y1, v, y2;
    for (i = 0; i < w * h; i += 2) {
        y1 = *src++;
        u = *src++;
        y2 = *src++;
        v = *src++;
        YUV444toRGB888(y1, u, v, dst);
        dst += 3;
        YUV444toRGB888(y2, u, v, dst);
        dst += 3;
    }
}

// Forward declaration of the Go callback function
extern void onFrameAvailableGo(void* data, int width, int height, int bytesPerPixel);

static LRESULT CALLBACK capVideoStreamCallback(HWND hwnd, LPVIDEOHDR vhdr) {
    fprintf(stderr, "Callback triggered\n");

    CAM *c = (CAM *)capGetUserData(hwnd);
    if (!c) {
        fprintf(stderr, "Failed to get user data\n");
        return 0;
    }

    fprintf(stderr, "Callback running\n");

    if (c->bmi.bmiHeader.biCompression == FOURCC('Y', 'U', 'Y', '2')) {
        YUYVToRGB24(c->w, c->h, (unsigned char*)vhdr->lpData, (unsigned char*)c->rgb);
    }

    onFrameAvailableGo(c->rgb, c->w, c->h, 3);
    return (LRESULT)TRUE;
}

// Manual capture
static BOOL captureFrame() {
	printf("Frame!\n");
    if (!capGrabFrameNoStop(_device->hwnd)) {
        fprintf(stderr, "Failed to grab frame\n");
        return FALSE;
    }
    return TRUE;
}

static int webcam_open(int deviceId, int width, int height) {
  // Allocate memory for the CAM structure
  _device = (CAM*) malloc(sizeof(CAM));
  if (!_device) {
    fprintf(stderr, "Failed to allocate memory for CAM structure\n");
    return -1;
  }

  // Create the capture window
  _device->hwnd = capCreateCaptureWindowA("Capture Window",     // Window name
                                          WS_CHILD,             // Window style
                                          0, 0, width, height,  // x, y, width, height
                                          GetDesktopWindow(),   // Parent window (set to desktop for testing)
                                          0);                   // Window ID

  if (!_device->hwnd) {
    fprintf(stderr, "Failed to create capture window\n");
    free(_device);
    return -1;
  }

  // Set the dimensions and allocate memory for RGB buffer
  _device->w = width;
  _device->h = height;
  _frame_size = width * height * 3;
  _device->rgb = malloc(_frame_size);
  if (!_device->rgb) {
    fprintf(stderr, "Failed to allocate memory for RGB buffer\n");
    DestroyWindow(_device->hwnd);
    free(_device);
    return -1;
  }

  // Connect the capture window to the device
  if (!capDriverConnect(_device->hwnd, deviceId)) {
    fprintf(stderr, "Failed to connect to capture device\n");
    DestroyWindow(_device->hwnd);
    free(_device->rgb);
    free(_device);
    return -1;
  }

  return 0;
}

static int webcam_start() {
    printf("Setting user data...\n");
    if (!capSetUserData(_device->hwnd, _device)) {
        fprintf(stderr, "Failed to set user data\n");
        return -1;
    }

    printf("Connecting capture driver...\n");
    if (!capDriverConnect(_device->hwnd, 0)) {
        fprintf(stderr, "Failed to connect to driver\n");
        return -1;
    }

    printf("Initializing BITMAPINFO structure...\n");
    memset(&_device->bmi, 0, sizeof(_device->bmi));
    _device->bmi.bmiHeader.biSize = sizeof(_device->bmi.bmiHeader);
    _device->bmi.bmiHeader.biWidth = 1280; // Supported width
    _device->bmi.bmiHeader.biHeight = 720; // Supported height
    _device->bmi.bmiHeader.biPlanes = 1;
    _device->bmi.bmiHeader.biBitCount = 16; // Supported bit count
    _device->bmi.bmiHeader.biCompression = 844715353; // 'YUY2' compression
    _device->bmi.bmiHeader.biSizeImage = 0;

    printf("Setting video format...\n");
    if (!capSetVideoFormat(_device->hwnd, &_device->bmi, sizeof(_device->bmi))) {
        fprintf(stderr, "Failed to set video format\n");
        return -1;
    }

    printf("Setting frame callback...\n");
    if (!capSetCallbackOnFrame(_device->hwnd, capVideoStreamCallback)) {
        fprintf(stderr, "Failed to set frame callback\n");
        return -1;
    }

    printf("Starting video stream...\n");
    if (!capCaptureSequenceNoFile(_device->hwnd)) {
       fprintf(stderr, "Failed to start video stream\n");
       return -1;
    }

    printf("Video stream started\n");
    return 0;
}

static int webcam_stop() {
    if (!capDriverDisconnect(_device->hwnd)) {
        fprintf(stderr, "Failed to disconnect driver\n");
        return -1;
    }
    return 0;
}

static void webcam_delete() {
    if (!_device) {
		return;
	}
	if (_device->hwnd) {
		DestroyWindow(_device->hwnd);
	}
	if (_device->rgb) {
		free(_device->rgb);
	}
	free(_device);
	_device = NULL;
}

static void postTestMessage() {
    PostMessage(_device->hwnd, WM_USER + 1, 0, 0); // Post a user-defined message
}

static HWND hwndTest;

// Function to run the message loop
static void runMessageLoop() {
    MSG msg;
	printf("Starting GetMessage()...\n");
	PostMessage(hwndTest, WM_USER + 1, 0, 0); // Post a user-defined message
    while (GetMessage(&msg, NULL, 0, 0)) {
        printf("GetMessage()!\n");
        TranslateMessage(&msg);
        DispatchMessage(&msg);

		// Manually capture frame for testing
		if (!captureFrame()) {
			printf("Failed to capture frame\n");
		}
    }
}

static void terminateMessageLoop() {
    PostQuitMessage(0);  // Posts a WM_QUIT message to the message queue
}

*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

var (
	loopIsRunning = false
)

//export onFrameAvailableGo
func onFrameAvailableGo(data unsafe.Pointer, width, height, bytesPerPixel C.int) {
	log.Println("Frame!")
	go func() {
		buf := make([]byte, width*height*bytesPerPixel)
		C.copyImage((*C.uint8_t)(unsafe.Pointer(&buf[0])), data)

		// Convert the buffer to an image.RGBA
		rgba := convertBGRAToRGBA(buf, int(width), int(height))

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
	if loopIsRunning {
		return fmt.Errorf("camera preview already running")
	}

	if C.webcam_start() != 0 {
		return fmt.Errorf("failed to start camera")
	}

	loopIsRunning = true
	go func() {
		C.runMessageLoop()
		loopIsRunning = false
	}()
	// Post a test message to ensure the message loop starts
	C.postTestMessage()
	return nil
}

func stopCamera() error {
	defer func() {
		if loopIsRunning {
			C.terminateMessageLoop()
			loopIsRunning = false
		}
	}()
	if C.webcam_stop() != 0 {
		return fmt.Errorf("failed to stop camera")
	}
	return nil
}

func closeCamera() {
	C.webcam_delete()
}
