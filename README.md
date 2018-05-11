imgurUpload

## Install

    go get -u github.com/weaming/imgurUpload

## Features

- upload with photo path
- upload with photo directory
- upload with photo url
- support annonymous upload

## Usage

1. Register app in [imgur settings](https://imgur.com/account/settings/apps)
1. Set app's redirect url to `http://127.0.0.1:1024` as callback to receive credentials
1. `expose IMGUR_CLIENT_ID=<your app client ID>`
1. `expose IMGUR_CLIENT_SECRET=<your app client secret>`

Command options

    Usage of imgurUpload:
      -a	anonymous mode will not upload to your album (default true)
      -p string
            target photo path/directory/url to upload
