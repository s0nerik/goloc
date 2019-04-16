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

## Supported platforms / formats

- Android
- iOS
- Flutter (experimental)
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

## License

Released under the [MIT License](https://github.com/s0nerik/goloc/blob/master/LICENSE).
