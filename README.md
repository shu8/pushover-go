# pushover-go 1.0.0

`pushover-go` is a simple command-line interface tool for sending messages to devices via [Pushover](https://pushover.net/). It is written in Go.

## Prerequisites

- A Pushover account -- see https://pushover.net/
- Pushover installed on your iOS, Android, or Desktop device
- A Pushover "Application Token" -- you need to register an application on Pushover's site to get this
- Your personal Pushover "User ID" -- you'll get this when you register with Pushover

## Installation

```
go get github.com/shu8/pushover-go
```

## Usage

```
pushover-go [options] 'Message to send'
echo 'Message to send' | pushover-go [options]
pushover-go [options] < file-with-message.txt

Options:
  -token string
        The application API token you want to send this message as. PUSHOVER_TOKEN env variable also available (required)
  -user string
        Your personal user ID to identify you as the sender. PUSHOVER_USER env variable also available. (required)
  -device string
        The name of the device you want to send the message to directly (optional)
  -priority int
        -2 no notification; -1 quiet notification; 1 high-priority; 2 require confirmation (optional)
  -sound string
        Name of the sound to play on the recipient device(s). See https://pushover.net/api#sounds (optional)
  -timestamp int
        Unix timestamp to show to the user (optional)
  -title string
        The title you want to give your message (optional)
  -url string
        A URL to show with your message (optional)
  -url-title string
        Text for the URL in --url to show (optional)
  -version
        Display the version of this tool and exit
```

## Examples

```
$ PUSHOVER_USER=user PUSHOVER_TOKEN=token pushover-go 'Hello World'
Sending message: Hello World
Successfully sent message
Total Quota: 7500
Quota Remaining: 7499
```

```
$ wget https://example.com/a-really-big-archive.zip && pushover-go -user user -token token 'Download complete'
Sending message: Download complete
Successfully sent message
Total Quota: 7500
Quota Remaining: 7498
```

```
$ cat <<EOT > my-message.txt
my long message
with many lines
EOT
$ pushover-go -user user -token token < my-message.txt
Sending message: my long message
with many lines

Successfully sent message
Total Quota: 7500
Quota Remaining: 7497
```

## License
[MIT License](./LICENSE)

## Notes

- There are a few other pushover CLI's out there, but I mainly wanted to make this to have a play around with Go!
- This tool supports [all of the options that Pushover accepts](https://pushover.net/api) which other ones might not (except for `attachment`, which I haven't needed to use yet!)
