# goloc

> A flexible tool for application localization using Google Sheets.

## Table of Contents

- [Features](#features)
- [Supported OS / architectures](#supported-os--architectures)
- [Supported platforms / formats](#supported-platforms--formats)
- [Setup](#setup)
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
- JSON

## Setup

- Download [latest release](https://github.com/s0nerik/goloc/releases/download/0.9/goloc.zip) and unpack it into your project's root folder
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

## Usage

- Create a script or build task definition with parameters best suited for your project. To see available parameters, run `goloc --help`. **goloc** is distributed in form of separate executables for each platform, so don't forget to take that into account when creating the localization script.
- Execute the script/task whenever you want to update localized strings. **goloc** will automatically replace any existing localization files with the updates ones.

## License

Released under the [MIT License](https://github.com/s0nerik/goloc/blob/master/LICENSE).
