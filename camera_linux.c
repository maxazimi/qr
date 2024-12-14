#include <assert.h>
#include <errno.h>
#include <fcntl.h>
#include <linux/videodev2.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/mman.h>
#include <sys/time.h>
#include <unistd.h>
#include "camera.h"

static const char DEVICE[] = "/dev/video0";

static int fd;
static struct {
    void *start;
    size_t length;
} *buffers;
static unsigned int num_buffers;
static struct v4l2_requestbuffers reqbuf = {0};

// External function declaration
extern void onFrameAvailableGo(void* data, int width, int height, int bytesPerPixel);

static int xioctl(int fd, int request, void *arg) {
    int r;
    do {
        r = ioctl(fd, request, arg);
    } while (-1 == r && EINTR == errno);
    return r;
}

static int init_mmap(void) {
    reqbuf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    reqbuf.memory = V4L2_MEMORY_MMAP;
    reqbuf.count = 5;
    if (-1 == xioctl(fd, VIDIOC_REQBUFS, &reqbuf)) {
        perror("VIDIOC_REQBUFS");
        return -1;
    }
    if (reqbuf.count < 2) {
        fprintf(stderr, "Not enough buffer memory\n");
        return -1;
    }
    buffers = calloc(reqbuf.count, sizeof(*buffers));
    if (!buffers) {
        fprintf(stderr, "Out of memory\n");
        return -1;
    }
    num_buffers = reqbuf.count;
    struct v4l2_buffer buffer;
    for (unsigned int i = 0; i < reqbuf.count; i++) {
        memset(&buffer, 0, sizeof(buffer));
        buffer.type = reqbuf.type;
        buffer.memory = V4L2_MEMORY_MMAP;
        buffer.index = i;
        if (-1 == xioctl(fd, VIDIOC_QUERYBUF, &buffer)) {
            perror("VIDIOC_QUERYBUF");
            return -1;
        }
        buffers[i].length = buffer.length;
        buffers[i].start = mmap(NULL, buffer.length, PROT_READ | PROT_WRITE, MAP_SHARED, fd, buffer.m.offset);
        if (MAP_FAILED == buffers[i].start) {
            perror("mmap");
            return -1;
        }
    }
    return 0;
}

static int init_device(int width, int height, int frame_rate) {
    struct v4l2_format fmt = {0};
    fmt.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    fmt.fmt.pix.width = width;
    fmt.fmt.pix.height = height;
    fmt.fmt.pix.pixelformat = V4L2_PIX_FMT_YUYV; // YUY2 format
    fmt.fmt.pix.field = V4L2_FIELD_NONE;
    if (-1 == xioctl(fd, VIDIOC_S_FMT, &fmt)) {
        perror("VIDIOC_S_FMT");
        return -1;
    }

    struct v4l2_streamparm param = {0};
    param.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    param.parm.capture.timeperframe.numerator = 1;
    param.parm.capture.timeperframe.denominator = frame_rate;
    if (-1 == xioctl(fd, VIDIOC_S_PARM, &param)) {
        perror("VIDIOC_S_PARM");
        return -1;
    }
    return init_mmap();
}

static int webcam_open(int width, int height, int frame_rate) {
    fd = open(DEVICE, O_RDWR);
    if (fd < 0) {
        perror(DEVICE);
        return -1;
    }
    return init_device(width, height, frame_rate);
}

static int webcam_start(void) {
    enum v4l2_buf_type type;
    struct v4l2_buffer buffer;
    for (unsigned int i = 0; i < num_buffers; i++) {
        memset(&buffer, 0, sizeof(buffer));
        buffer.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        buffer.memory = V4L2_MEMORY_MMAP;
        buffer.index = i;
        if (-1 == xioctl(fd, VIDIOC_QBUF, &buffer)) {
            perror("VIDIOC_QBUF");
            return -1;
        }
    }
    type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    if (-1 == xioctl(fd, VIDIOC_STREAMON, &type)) {
        perror("VIDIOC_STREAMON");
        return -1;
    }
    return 0;
}

static int webcam_stop(void) {
    enum v4l2_buf_type type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    if (-1 == xioctl(fd, VIDIOC_STREAMOFF, &type)) {
        perror("VIDIOC_STREAMOFF");
        return -1;
    }
    return 0;
}

static void webcam_close(void) {
    for (unsigned int i = 0; i < num_buffers; i++) {
        munmap(buffers[i].start, buffers[i].length);
    }
    free(buffers);
    close(fd);
}

static int capture_frame(void) {
    struct v4l2_buffer buffer;
    memset(&buffer, 0, sizeof(buffer));
    buffer.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    buffer.memory = V4L2_MEMORY_MMAP;
    if (-1 == xioctl(fd, VIDIOC_DQBUF, &buffer)) {
        switch (errno) {
        case EAGAIN:
            return 0;
        case EIO:
        default:
            perror("VIDIOC_DQBUF");
            return -1;
        }
    }
    assert(buffer.index < num_buffers);

    // Allocate memory for RGB data
    uint8_t *rgb_data = malloc(640 * 480 * 3);
    if (!rgb_data) {
        fprintf(stderr, "Out of memory\n");
        return -1;
    }

    // Convert YUY2 to RGB24
    YUY2toRGB24((uint8_t*)buffers[buffer.index].start, rgb_data, 640, 480);

    // Call the Go callback function with RGB data
    onFrameAvailableGo(rgb_data, 640, 480, 3);

    // Free the allocated RGB data memory
    free(rgb_data);

    if (-1 == xioctl(fd, VIDIOC_QBUF, &buffer)) {
        perror("VIDIOC_QBUF");
        return -1;
    }
    return 1;
}
