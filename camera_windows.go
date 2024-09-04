package camera

/*
#cgo LDFLAGS: -lvfw32
#include <windows.h>
#include <vfw.h>

#define FOURCC(a, b, c, d) ((DWORD)(((d) << 24) | ((c) << 16) | ((b) << 8) | (a)))

typedef struct {
	HWND hwnd;
	long w, h;
	BITMAPINFO bmi;
	void *rgb;
} CAM;

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
	for (i = 0; i < w * h; i += 2)
	{
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

LRESULT CALLBACK capVideoStreamCallback(HWND hwnd, LPVIDEOHDR vhdr) {
	CAM *c = (CAM *)capGetUserData(hwnd);
	if (c->bmi.bmiHeader.biCompression == FOURCC('Y', 'U', 'Y', '2')) {
		YUYVToRGB24(c->w, c->h, vhdr->lpData, c->rgb);
	}

	onFrameAvailableGo(c->rgb, c->w, c->h, 3);
	return 0;
}

// Global device pointer
static unsigned long _frame_size = 0;
static CAM* _device = NULL;

static void grabFrame(uint8_t *buffer) {
	capGrabFrameNoStop(_device->hwnd);
	memcpy(buffer, _device->rgb, _frame_size);
}

static int webcam_open(int deviceId, int width, int height) {
	_device = (CAM*)malloc(sizeof(CAM));
	_device->hwnd = capCreateCaptureWindow(0, 0, 0, 0, 0, 0, 0, 0);
	_device->w = width;
	_device->h = height;

	_frame_size = width * height * 3;
	_device->rgb = malloc(_frame_size);
	return 0;
}

static int webcam_start() {
	capSetUserData(_device->hwnd, _device);
	capDriverConnect(_device->hwnd, 0);
	capGetVideoFormat(_device->hwnd, &_device->bmi, sizeof(_device->bmi));

	_device->bmi.bmiHeader.biWidth = _device->w;
	_device->bmi.bmiHeader.biHeight = _device->h;
	_device->bmi.bmiHeader.biSizeImage = 0;

	capSetVideoFormat(_device->hwnd, &_device->bmi, sizeof(_device->bmi));
	capSetCallbackOnFrame(_device->hwnd, capVideoStreamCallback);

	return 0;
}

static int webcam_stop() {
	capDriverDisconnect(_device->hwnd);
	return 0;
}

static void webcam_delete() {
	DestroyWindow(_device->hwnd);
	free(_device->rgb);
	free(_device);
	_device = NULL;
}

*/
import "C"
import (
	"fmt"
	"image"
	"unsafe"
)

var (
	frameBufferChan = make(chan *image.RGBA, 1)
)

//export onFrameAvailableGo
func onFrameAvailableGo(data unsafe.Pointer, width, height, bytesPerPixel C.int) {
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

func OpenCamera(id, width, height int) error {
	frameBufferChan = make(chan *image.RGBA, 10)

	if C.webcam_open(C.int(id), C.int(width), C.int(height)) != 0 {
		return fmt.Errorf("failed to initialize camera")
	}
	return nil
}

func StartCamera() error {
	if C.webcam_start() != 0 {
		return fmt.Errorf("failed to start camera")
	}
	return nil
}

func StopCamera() error {
	if C.webcam_stop() != 0 {
		return fmt.Errorf("failed to stop camera")
	}
	return nil
}

func CloseCamera() {
	C.webcam_delete()
}

func GetCameraFrameChan() <-chan *image.RGBA {
	return frameBufferChan
}
