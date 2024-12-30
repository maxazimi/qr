//go:build android

package platform

//#include <jni.h>
import "C"
import (
	"gioui.org/app"
	_ "gioui.org/app/permission/camera"
	"gioui.org/io/event"
	"git.wow.st/gmp/jni"
	"github.com/maxazimi/qr/ui/ev"
	"log"
	"unsafe"
)

const (
	cameraPermission = "android.permission.CAMERA"
)

var (
	viewEvent app.ViewEvent
)

func HandleEvent(e event.Event) {
	switch e := e.(type) {
	case app.ViewEvent:
		viewEvent = e
	default:
	}
}

func RequestCameraPermission() chan error {
	permission := cameraPermission
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(errChan)
		}()

		errChan <- requestPermission(permission)
		return

		jvm := jni.JVMFor(app.JavaVM())
		err := jni.Do(jvm, func(env jni.Env) error {
			uptr := app.AppContext()
			appCtx := *(*jni.Object)(unsafe.Pointer(&uptr))

			// Get the class and method IDs for checkSelfPermission
			activityCls := jni.FindClass(env, "android/app/Activity")
			checkPermissionMID := jni.GetMethodID(env, activityCls, "checkSelfPermission", "(Ljava/lang/String;)I")

			jPermission := jni.JavaString(env, permission)

			// Check if camera permission is granted
			permissionGranted, err := jni.CallIntMethod(env, appCtx, checkPermissionMID, jni.Value(jPermission))
			if err != nil {
				return err
			}
			if permissionGranted == 0 {
				return nil
			}
			//return fmt.Errorf("camera permission not granted")

			// Prepare the permissions array
			stringClass := jni.FindClass(env, "java/lang/String")
			jPermissions := jni.NewObjectArray(env, 1, stringClass, jni.Object(jPermission))
			err = jni.SetObjectArrayElement(env, jPermissions, 0, jni.Object(jPermission))
			if err != nil {
				return err
			}

			// Call the requestPermissions method
			requestPermissionMID := jni.GetMethodID(env, activityCls, "requestPermissions", "([Ljava/lang/String;I)V")
			return jni.CallVoidMethod(env, appCtx, requestPermissionMID, jni.Value(jPermissions), jni.Value(1))
		})
		if err != nil {
			log.Printf("requestPermission() JVM error: %s", err)
		}
	}()
	return errChan
}

func requestPermission(permission string) error {
	const name = "requestPermission"
	jvm := jni.JVMFor(app.JavaVM())
	return jni.Do(jvm, func(env jni.Env) error {
		ptr := app.AppContext()
		context := *(*jni.Object)(unsafe.Pointer(&ptr))

		class := jni.GetObjectClass(env, context)
		methodID := jni.GetMethodID(env, class, name, "(Ljava/lang/String;)V")

		param := jni.Value(jni.JavaString(env, permission))
		return jni.CallVoidMethod(env, context, methodID, param)
	})
}

func OpenURL(url string) error {
	const name = "openUrl"
	jvm := jni.JVMFor(app.JavaVM())
	return jni.Do(jvm, func(env jni.Env) error {
		ptr := app.AppContext()
		context := *(*jni.Object)(unsafe.Pointer(&ptr))

		class := jni.GetObjectClass(env, context)
		methodID := jni.GetMethodID(env, class, name, "(Ljava/lang/String;)V")

		param := jni.Value(jni.JavaString(env, url))
		return jni.CallVoidMethod(env, context, methodID, param)
	})
}

func ShowText(text string) {
	go func() {
		jvm := jni.JVMFor(app.JavaVM())
		err := jni.Do(jvm, func(env jni.Env) error {
			uptr := app.AppContext()
			appCtx := *(*jni.Object)(unsafe.Pointer(&uptr))

			class := jni.GetObjectClass(env, appCtx)
			methodID := jni.GetMethodID(env, class, "showText", "(Ljava/lang/String;)V")

			param := jni.Value(jni.JavaString(env, text))
			return jni.CallVoidMethod(env, appCtx, methodID, param)
		})
		if err != nil {
			log.Printf("ShowText() jvm error: %s", err)
		}
	}()
}

//export Java_com_maxazimi_qr_App_onActivityCreated
func Java_com_maxazimi_qr_App_onActivityCreated(env *C.JNIEnv, cls C.jclass, width, height C.int) {
	ev.RequestEvent(ev.AppLoadEvent{})
	ev.RequestEvent(ev.WindowSizeEvent{
		Width:  int(width),
		Height: int(height),
	})
}
