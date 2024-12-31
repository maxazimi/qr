# QR Reader

QR Reader is a cross-platform QR code reader built with the GioUI graphics library. It supports macOS, Linux, Windows, and Android.

## Features
- Cross-platform compatibility
- Built with GioUI graphics library
- Support for macOS, Linux, Windows, and Android

## Dependencies
- **macOS**: Requires Xcode Command Line Tools.
- **Linux**: Requires v4l2 (Video for Linux 2).
- **Windows**: Requires MinGW or MSYS2.
- **Android**: Requires NDK 24 (Native Development Kit).

## Camera Package
This app uses the [Camera Package](https://github.com/maxazimi/camera) for accessing the camera across different platforms. The `camera` package is a cross-platform solution written in CGo, supporting macOS, Linux, Windows, and Android.

## Usage
Build for desktop (macOS, Linux, Windows)
```shell
make
```
Build for Android
```shell
make android
make apk_install
```
Modify the Makefile before build.

## Contributing
Contributions are highly appreciated.

## TODO
- [ ] **Fix Bugs**

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits
- QR code reading powered by [gozxing](https://github.com/makiuchi-d/gozxing/qrcode)
- GUI powered by [GioUI](https://gioui.org)

## Contact

For any questions, suggestions, or contributions, please contact:
- **Email**: [maxazimy@gmail.com](mailto:maxazimy@gmail.com)
