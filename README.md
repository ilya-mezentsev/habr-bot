# Habr-bot

Bot for parsing articles from [Habr](https://habr.com)

## Deployment
```bash
$ make workspace && make build
```

## Testing
```bash
$ make test
```

## Debug
```bash
$ make run-cli
```

## Telegram

To use Telegram mode
* Add bot token ```TG_BOT_TOKEN=<your_bot_token>``` (see [this link](https://core.telegram.org/bots)) in ```.env``` file
* Run
```bash
$ make run-tg
```
