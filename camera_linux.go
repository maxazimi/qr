//go:build linux && !android && !ios

package camera

/*
#cgo LDFLAGS: -lv4l2

// https://github.com/kosua20/sr_webcam/blob/master/src/sr_webcam_lin.c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/ioctl.h>
#include <asm/types.h>
#include <linux/videodev2.h>
#include <pthread.h>
#include <math.h>
#include "camera.h"

typedef struct {
	void* start;
	long length;
} _webcam_buffer;

typedef struct {
	int fid;
	webcam_device_t* parent;
	int width;
	int height;
	int id;
	int framerate;
	_webcam_buffer* buffers;
	int buffersCount;
	pthread_t thread;
} _webcam_v4lInfos;

int _webcam_wait_ioctl(int fid, int request, void* arg) {
	int r;
	do {
		r = ioctl(fid, request, arg);
	} while (r == -1 && EINTR == errno);
	return r;
}

void* _webcam_callback_loop(void* arg) {
	_webcam_v4lInfos* stream = (_webcam_v4lInfos*)arg;

	// TODO Make sure that this is blocking to avoid overload.
	while (1) {
		fd_set fds;
		FD_ZERO(&fds);
		FD_SET(stream->fid, &fds);

		struct timeval tv;
		tv.tv_sec  = 2;
		tv.tv_usec = 0;

		int res = select(stream->fid + 1, &fds, NULL, NULL, &tv);
		if (res == -1 || res == 0) {
			return NULL;
		}

		struct v4l2_buffer buf;
		memset(&buf, 0, sizeof(buf));
		buf.type   = V4L2_BUF_TYPE_VIDEO_CAPTURE;
		buf.memory = V4L2_MEMORY_MMAP;

		if (_webcam_wait_ioctl(stream->fid, VIDIOC_DQBUF, &buf) == -1) {
			if (errno != EIO) {
				return NULL;
			}
		}

		unsigned int dstSize = stream->width * stream->height * 3;
		unsigned char* data	 = (unsigned char*)malloc(dstSize);
		const char* srcData	 = (const char*)(stream->buffers[buf.index].start);
		for (int y = 0; y < stream->height; ++y) {
			for (int x = 0; x < stream->width; ++x) {
				const int baseShift		= (y * stream->width + x);
				data[3 * baseShift + 0] = srcData[4 * baseShift + 0];
				data[3 * baseShift + 1] = srcData[4 * baseShift + 1];
				data[3 * baseShift + 2] = srcData[4 * baseShift + 2];
			}
		}
		stream->parent->callback(stream->parent, data);
		_webcam_wait_ioctl(stream->fid, VIDIOC_QBUF, &buf);
	}
	return NULL;
}

// Global device pointer
static webcam_device_t* device = NULL;

static inline uint32_t webcam_format_size() {
	return (uint32_t)(device->width) * (uint32_t)(device->height) * 4;
}

static void copyImage(uint8_t *dstBuf, void* srcBuf) {
	memcpy(dstBuf, srcBuf, webcam_format_size());
}

int webcam_open(int deviceId, int width, int height, int framerate) {
	if (device != NULL) {
		return -1;
	}

	device = (webcam_device_t*)malloc(sizeof(webcam_device_t));
	if (!device) {
		return -1;
	}

	memset(device, 0, sizeof(*device));
	device->deviceId  = deviceId;
	device->width	  = width;
	device->height	  = height;
	device->framerate = framerate;

	// Already setup.
	if (device->stream) {
		return -1;
	}
	if (device->deviceId < 0) {
		return -1;
	}

	_webcam_v4lInfos* stream = (_webcam_v4lInfos*)malloc(sizeof(_webcam_v4lInfos));
	memset(stream, 0, sizeof(_webcam_v4lInfos));
	stream->parent = device;
	stream->fid	   = -1;

	// Try to open the corresponding handle.
	for (int i = device->deviceId; i >= 0; --i) {
		char file[256];
		snprintf(file, 255, "/dev/video%d", i);
		int fid = open(file, O_RDONLY);
		if (fid >= 0) {
			// We found a valid file.
			stream->id = i;
			close(fid);

			// Open the file for real.
			stream->fid = open(file, O_RDWR | O_NONBLOCK, 0);
			break;
		}
	}

	// Failed to find any device.
	if (stream->fid < 0) {
		free(stream);
		return -1;
	}
	int fid = stream->fid;

	// Configure the device.
	struct v4l2_capability cap;
	// If we can't query the device capabilities, or it doesn't support video streaming, skip.
	if (_webcam_wait_ioctl(fid, VIDIOC_QUERYCAP, &cap) == -1 || !(cap.capabilities & V4L2_CAP_VIDEO_CAPTURE) || !(cap.capabilities & V4L2_CAP_STREAMING)) {
		free(stream);
		return -1;
	}

	// Select output crop.
	struct v4l2_cropcap cropCap;
	cropCap.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
	_webcam_wait_ioctl(fid, VIDIOC_CROPCAP, &cropCap);
	struct v4l2_crop crop;
	crop.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
	crop.c	  = cropCap.defrect; // Default rectangle.
	_webcam_wait_ioctl(fid, VIDIOC_S_CROP, &crop);

	// Select output format.
	struct v4l2_format fmt;
	memset(&fmt, 0, sizeof(fmt));
	fmt.type		   = V4L2_BUF_TYPE_VIDEO_CAPTURE;
	fmt.fmt.pix.width  = device->width;
	fmt.fmt.pix.height = device->height;
	// Force RGB.
	fmt.fmt.pix.pixelformat = V4L2_PIX_FMT_RGB32;
	// Pixel format not supported.
	if (_webcam_wait_ioctl(fid, VIDIOC_S_FMT, &fmt) == -1) {
		free(stream);
		return -1;
	}
	fmt.fmt.pix.bytesperline = fmax(fmt.fmt.pix.bytesperline, fmt.fmt.pix.width * 2);
	fmt.fmt.pix.sizeimage	 = fmax(fmt.fmt.pix.sizeimage, fmt.fmt.pix.bytesperline * fmt.fmt.pix.height);

	// Update the size based on the format constraints.
	stream->width  = fmt.fmt.pix.width;
	stream->height = fmt.fmt.pix.height;

	// Allocate buffers for video frames.
	struct v4l2_requestbuffers req;
	memset(&req, 0, sizeof(req));
	req.count  = 4;
	req.type   = V4L2_BUF_TYPE_VIDEO_CAPTURE;
	req.memory = V4L2_MEMORY_MMAP;

	// If we can't get at least two buffers, skip.
	if (_webcam_wait_ioctl(fid, VIDIOC_REQBUFS, &req) == -1 || req.count < 2) {
		free(stream);
		return -1;
	}
	_webcam_buffer* buffers = calloc(req.count, sizeof(_webcam_buffer));
	if (!buffers) {
		free(stream);
		return -1;
	}
	stream->buffersCount = req.count;

	// Allocate the buffers.
	for (int bid = 0; bid < (int)(req.count); ++bid) {
		struct v4l2_buffer buf;
		memset(&buf, 0, sizeof(buf));
		buf.type   = V4L2_BUF_TYPE_VIDEO_CAPTURE;
		buf.memory = V4L2_MEMORY_MMAP;
		buf.index  = bid;

		if (_webcam_wait_ioctl(fid, VIDIOC_QUERYBUF, &buf) == -1) {
			for (int obid = 0; obid < bid; ++obid) {
				munmap(buffers[obid].start, buffers[obid].length);
			}
			free(buffers);
			free(stream);
			return -2;
		}

		buffers[bid].length = buf.length;
		buffers[bid].start	= mmap(NULL, buf.length, PROT_READ | PROT_WRITE, MAP_SHARED, fid, buf.m.offset);
		if (buffers[bid].start == MAP_FAILED) {
			for (int obid = 0; obid < bid; ++obid) {
				munmap(buffers[obid].start, buffers[obid].length);
			}
			free(buffers);
			free(stream);
			return -2;
		}
	}
	stream->buffers = buffers;

	// Try to set the framerate.
	struct v4l2_streamparm fpsParams;
	memset(&fpsParams, 0, sizeof(fpsParams));
	fpsParams.type									= V4L2_BUF_TYPE_VIDEO_CAPTURE;
	fpsParams.parm.capture.timeperframe.numerator	= 1;
	fpsParams.parm.capture.timeperframe.denominator = (unsigned int)(device->framerate);
	_webcam_wait_ioctl(fid, VIDIOC_S_PARM, &fpsParams);
	_webcam_wait_ioctl(fid, VIDIOC_G_PARM, &fpsParams);
	stream->framerate = (int)fpsParams.parm.capture.timeperframe.denominator;

	// Update the device infos.
	device->stream	  = stream;
	device->width	  = stream->width;
	device->height	  = stream->height;
	device->deviceId  = stream->id;
	device->framerate = stream->framerate;

	return 0;
}

int webcam_start() {
	if (device == NULL) {
		return -1;
	}

	if (device->stream && device->running == 0) {
		_webcam_v4lInfos* stream = (_webcam_v4lInfos*)(device->stream);

		// Prepare all buffers.
		for (int bid = 0; bid < stream->buffersCount; ++bid) {
			struct v4l2_buffer buf;
			memset(&buf, 0, sizeof(buf));
			buf.type   = V4L2_BUF_TYPE_VIDEO_CAPTURE;
			buf.memory = V4L2_MEMORY_MMAP;
			buf.index  = bid;
			if (_webcam_wait_ioctl(stream->fid, VIDIOC_QBUF, &buf) == -1) {
				return -1;
			}
		}

		enum v4l2_buf_type type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
		if (_webcam_wait_ioctl(stream->fid, VIDIOC_STREAMON, &type) == -1) {
			return -1;
		}
		pthread_create(&stream->thread, NULL, &_webcam_callback_loop, device->stream);
		device->running = 1;
	}
	return 0;
}

int webcam_stop() {
	if (device->stream && device->running == 1) {
		_webcam_v4lInfos* stream = (_webcam_v4lInfos*)(device->stream);
		enum v4l2_buf_type type	 = V4L2_BUF_TYPE_VIDEO_CAPTURE;
		if (_webcam_wait_ioctl(stream->fid, VIDIOC_STREAMOFF, &type) == -1) {
			return -1;
		}

		pthread_cancel(stream->thread);
		device->running = 0;
	}
	return 0;
}

void webcam_delete() {
	if (device == NULL) {
		return;
	}
	if (device->running == 1) {
		webcam_stop(device);
	}

	if (device->stream) {
		// Unmap and delete all buffers.
		_webcam_v4lInfos* stream = (_webcam_v4lInfos*)(device->stream);
		_webcam_buffer* buffers	= stream->buffers;
		for (int bid = 0; bid < stream->buffersCount; ++bid) {
			munmap(buffers[bid].start, buffers[bid].length);
		}
		free(buffers);
		close(stream->fid);
	}

	free(device);
	device = NULL;
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
func onFrameAvailableGo(data unsafe.Pointer, w, h C.int) {
	go func() {
		buf := make([]byte, C.webcam_format_size())
		C.copyImage((*C.uint8_t)(unsafe.Pointer(&buf[0])), data)

		// Convert the buffer to an image.RGBA
		rgba := convertBGRAToRGBA(buf, int(w), int(h))

		select {
		case frameBufferChan <- rgba:
		default:
			// Drop the frame if the channel is full
		}
	}()
}

func OpenCamera(id, width, height int) error {
	frameBufferChan = make(chan *image.RGBA, 10)

	if C.webcam_open(C.int(id), C.int(width), C.int(height), 30) != 0 {
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
