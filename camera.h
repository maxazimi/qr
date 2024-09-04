#ifndef CAMERA_H
#define CAMERA_H

struct _webcam_device;
typedef struct _webcam_device webcam_device_t;
typedef void (*webcam_callback_t)(webcam_device_t* device, void* data);

struct _webcam_device {
	int deviceId;
	int width;
	int height;
	int framerate;
	int running;

	void* stream;
	webcam_callback_t callback;
};

#if defined(__OBJC__)
#import <Foundation/Foundation.h>
#import <AVFoundation/AVFoundation.h>

@interface WebcamVideoStream : NSObject <AVCaptureVideoDataOutputSampleBufferDelegate> {
	AVCaptureSession* _captureSession;
	AVCaptureDevice* _captureDevice;
	AVCaptureVideoDataOutput* _captureDataOut;
	AVCaptureDeviceInput* _captureDataIn;

@public
	webcam_device_t* _parent;
@public
	int _width;
@public
	int _height;
@public
	int _id;
@public
	int _framerate;
}

- (BOOL)setupWithID:(int)deviceID rate:(int)framerate width:(int)w height:(int)h;
- (void)start;
- (void)stop;

@end

#endif // defined(__OBJC__)
#endif // CAMERA_H
