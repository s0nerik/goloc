# goloc

> A flexible tool for application localization using Google Sheets.

## Table of Contents

- [Features](#features)
- [Supported platforms / formats](#supported-platforms--formats)
- [Setup](#setup)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Easy configuration
- High configurability
- Precise error reporting
- Multiple supported target platforms
- Customizable format strings
- Missing localization reports

## Supported platforms / formats

- Android
- iOS
- JSON

## Setup

- Download latest release and unpack it into your project's root folder
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
- Create a script or build task definition with parameters best suited for the project

## Usage

TODO

## Contributing

TODO

## License

Released under the [MIT License](https://github.com/s0nerik/goloc/blob/master/LICENSE).
