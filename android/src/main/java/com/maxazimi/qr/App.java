package com.maxazimi.qr;

import android.app.Activity;
import android.app.Application;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.net.Uri;
import android.os.Bundle;
import android.util.DisplayMetrics;
import android.widget.Toast;

import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import org.gioui.Gio;

public class App extends Application implements Application.ActivityLifecycleCallbacks {
	private static final int PERMISSION_REQUEST = 1;
	private Activity mCurrentActivity = null;

	public Activity getCurrentActivity() {
		return mCurrentActivity;
	}

	@Override
	public void onActivityCreated(Activity activity, Bundle bundle) {
		mCurrentActivity = activity;
		DisplayMetrics displayMetrics = new DisplayMetrics();
		mCurrentActivity.getWindowManager().getDefaultDisplay().getMetrics(displayMetrics);

		//Intent browserIntent = new Intent(Intent.ACTION_VIEW, Uri.parse("http://www.google.com"));
		//startActivity(browserIntent);
	}
	@Override
	public void onActivityStarted(Activity activity) {
	}
	@Override
	public void onActivityResumed(Activity activity) {
	}
	@Override
	public void onActivityPaused(Activity activity) {
	}
	@Override
	public void onActivityStopped(Activity activity) {
	}
	@Override
	public void onActivitySaveInstanceState(Activity activity, Bundle bundle) {}
	@Override
	public void onActivityDestroyed(Activity activity) {
	}

	@Override
	public void onCreate() {
		super.onCreate();
		Gio.init(this);
		registerActivityLifecycleCallbacks(this);
	}

	public void showText(String text) {
		if (mCurrentActivity == null) return;
		mCurrentActivity.runOnUiThread(new Runnable() {
			public void run() {
				Toast.makeText(mCurrentActivity, text, Toast.LENGTH_SHORT).show();
			}
		});
	}

	private void requestPermission(String permission) {
		if (ContextCompat.checkSelfPermission(this, permission)
				!= PackageManager.PERMISSION_GRANTED) {
			ActivityCompat.requestPermissions(
					mCurrentActivity,
					new String[] { permission },
					PERMISSION_REQUEST);
		}
	}
}
