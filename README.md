# Habr-bot

Bot for parsing articles from habr.com

## Deployment
```bash
$ git clone https://github.com/ilya-mezentsev/habr_bot.git

$ cd habr_bot && bash prepare_workspace.sh $(pwd)
```

## Testing
```bash
$ bash run.sh go_tests
```

## Debug
```bash
$ bash run.sh cli
```

## Telegram

To use Telegram mode
* Add bot token ```TG_BOT_TOKEN=<your_bot_token>``` (see [this link](https://core.telegram.org/bots)) in ```.env``` file
* Run
```bash
$ bash run.sh tg
```
