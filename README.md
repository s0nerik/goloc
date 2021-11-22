![Logo](docs/images/goloc.png?raw=true)

# goloc

> A flexible tool for application localization using Google Sheets.

## Table of Contents

- [Features](#features)
- [Supported OS / architectures](#supported-os--architectures)
- [Supported platforms / formats](#supported-platforms--formats)
- [Setup](#setup)
- [Localization document](#localization-document)
	- [Localizations sheet](#localizations-sheet)
	- [Formats sheet](#formats-sheet)
- [Usage](#usage)
	- [Android](#android)
	- [Flutter](#flutter)
- [macOS Catalina usage notes](#macos-catalina-usage-notes)
- [License](#license)

## Features

- Easy configuration
- High configurability
- Precise error reporting
- Multiple supported target platforms
- Customizable format strings
- Missing localization reports

## Supported OS / architectures

**goloc** can be built for each OS/architecture supported by golang, but release archives
contain binaries only for **amd64** architecture for **macOS**, **Linux** and **Windows**.

## Supported formats

- [Android](#android)
- iOS
- [Flutter](#flutter)
- JSON

## Setup

- Download a `goloc.zip` file from the [latest release](https://github.com/s0nerik/goloc/releases/latest) and unpack it into your project's root folder
- Download `client_secret.json` file from Google API Console and put it inside a newly created `goloc` folder. To do so, follow these steps:
	- Open [Google API Console](https://console.developers.google.com)
	- Select a project (or create a new one)
	- Press `ENABLE APIS AND SERVICES` button
	- Find `Google Sheets API`
	- Press `ENABLE`
	- Go to `Dashboard->Credentials`
	- Press `Create credentials->Service account key`
	- Select any type of service account and `JSON` for key type, then press `Create`
	- Rename the downloaded file to `client_secret.json` and put it into a `goloc` folder of a project
- Create a new [localization document](#localization-document)
- Share your localization document with a service account created previously. To do so, follow these steps:
	- Open the `client_secret.json` file
	- Copy the `client_email` value
	- Open the localization document
	- Press `SHARE` button
	- Paste the `client_email` value into the `People` input field.

## Localization document

Each localization document consists of **formats** and **localizations** sheets. One localization document can have multiple sheets for both.

The simplest way to create a new **goloc**-compatible localization document is to copy the [sample spreadsheet](https://docs.google.com/spreadsheets/d/1pmPPYLrHfSGLM-1MPYEGtbb9Z5iHFUL-xqXNFS0DyaM/edit?usp=sharing). However, you can easily create a **goloc**-compatible localization document yourself just by following the simple requirements described below.

### Localizations sheet

![Example localizations sheet](docs/images/localizations_example.jpg?raw=true)

On the example above you can see a **goloc**-compatible localizations sheet. The rules to make a localizations sheet **goloc**-compatible are:

- First row must contain column names
- There must be exactly one **key** column and at least one **language** column
- **Key** column can have any name, but the dafault name is `key`
- Each **language** column must be named as `lang_<lanaguage code>`
- To define a format string, you can use `{format_name}` in place of the formatted value (each format must be specified in the [formats sheet](#formats-sheet))

### Formats sheet

![Example formats sheet](docs/images/formats_example.jpg?raw=true)

On the example above you can see a **goloc**-compatible formats sheet. The rules to make a formats sheet **goloc**-compatible are:

- First row must contain column names
- There must be exactly one **format** column and at least one **platform** column
- **Format** column can have any name, but the dafault name is `format`
- Each **platform** column must have a name of a [**goloc**-supported platform](https://github.com/s0nerik/goloc/tree/master/platforms).
- Empty format name can be used to define a default format (used as `{}`)

## Usage

- Create a script or build task definition with parameters best suited for your project. To see available parameters, run `goloc --help`. **goloc** is distributed in form of separate executables for each platform, so don't forget to take that into account creating localization script.
- Execute the script/task whenever you want to update localized strings. **goloc** will automatically replace any existing localization files with the updated ones.

### Android

No special configuration in code is required.

Example **gradle** task specification:

```gradle
task "fetchLocalizations"(type: Exec) {
    def osName = System.getProperty('os.name').toLowerCase()
    def isWindows = osName.contains("win")
    def isMac = osName.contains("mac")
    def isUnix = osName.contains("nix") || osName.contains("nux") || osName.contains("aix")

    def params = [
            '--credentials', "goloc/client_secret.json",
            '--platform', 'android',
            '--spreadsheet', '1MbtglvGyEey3gH8yh4c9QovCIbtl5EcwqWqTZUiNga8',
            '--tab', "localizations",
            "--key-column", "key",
            '--resources', "app/src/main/res/",
            '--default-localization', 'en',
            '--default-localization-file-path', "app/src/main/res/values/localized_strings.xml"
    ]

    if (isWindows) {
        params = ['cmd', '/c', 'goloc\\windows_amd64.exe'] + params
    } else if (isMac) {
        params = ['./goloc/darwin_amd64'] + params
    } else if (isUnix) {
        params = ['./goloc/linux_amd64'] + params
    } else {
        logger.error('Your OS is not supported.')
        return
    }

    commandLine params
}
```

### Flutter

Localized strings can be accessed through `AppLocalizations.of(context)`

Requirements:

- Add `sprintf: ^6.0.0` to the `dependencies` section of `pubspec.yaml`
- Add `AppLocalizationsDelegate()` to `localizationsDelegates` of the app widget constructor
- Specify supported localizations in `supportedLocales` of the app widget constructor
- (Recommended) Add `DefaultIntlLocaleDelegate()` to `localizationsDelegates` of the app widget constructor. This will make `intl`-dependent formatters use currently selected locale.

```dart
class DefaultIntlLocaleDelegate extends LocalizationsDelegate<Null> {
  @override
  bool isSupported(Locale locale) => true;

  @override
  Future<Null> load(Locale locale) {
    Intl.defaultLocale = locale.toLanguageTag();
    return Future.value(null);
  }

  @override
  bool shouldReload(LocalizationsDelegate<AppLocalizations> old) => false;
}
```

Example **bash** localization script:

```bash
#!/bin/bash

case "$OSTYPE" in
  darwin*)  EXECUTABLE="darwin_amd64" ;;
  linux*)   EXECUTABLE="linux_amd64" ;;
  msys*)    EXECUTABLE="windows_amd64.exe" ;;
  *)
	  echo "Platform is not supported: $OSTYPE"
	  exit 1
  ;;
esac

goloc/${EXECUTABLE} -c goloc/client_secret.json -p flutter -s 1MbtglvGyEey3gH8yh4c9QovCIbtl5EcwqWqTZUiNga8 -t localizations -r lib/intl
```

## macOS Catalina usage notes

Due to the security improvements in the macOS Catalina, any 3rd party application downloaded from the internet has to be notarized to be launched without additional actions from the user side. Since **goloc** is entirely free, I can't afford Apple Developer Program subscription for notarizing macOS builds. Luckily, Apple has left a way to launch a non-notarized app, but it requires some actions.

Here's the instruction on how to launch **goloc** on macOS Catalina:

1. Upon the first launch you'll see a window like this:![catalina_0](docs/images/catalina_0.png?raw=true)
2. Go to settings and choose `Security & Privacy`:![catalina_1](docs/images/catalina_1.png?raw=true)
3. Choose `Open Anyway`:![catalina_2](docs/images/catalina_2.png?raw=true)
4. Re-launch **goloc**, and you'll see a next window:![catalina_3](docs/images/catalina_3.png?raw=true)
5. Choose `Open`. Next time you launch **goloc**, macOS won't complain anymore.

## License

Released under the [MIT License](https://github.com/s0nerik/goloc/blob/master/LICENSE).
