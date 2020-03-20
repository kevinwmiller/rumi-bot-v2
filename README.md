# Requirements
- Go 1.12 or 13 or something. It doesn't matter. Just have modules

# Setup

- Create .env with the following variables
```
DISCORD_BOT_TOKEN=
DISCORD_NOTIFICATION_CHANNEL_ID=
TWITCH_CLIENT_ID=
TWITCH_WATCH_USERNAME=
```

## DISCORD_BOT_TOKEN
- Go to https://discordapp.com/developers/applications and create a new application.
- Click the bot tab on the left
- Under Token, click Copy
- This is your bot token

## DISCORD_NOTIFICATION_CHANNEL_ID
- Go to your user settings, and then Appearance
- Scroll to the bottom and turn on Developer Mode
    - This will give you additional context menu items when you right click a channel
- Find the text channel you want to notify when the stream starts/stops, right click it, and at the bottom, click Copy ID
- This is your notification channel ID

## TWITCH_CLIENT_ID
- Go to https://dev.twitch.tv/console
- Click Register Your Application
- Give it a name
- Add some URL for the OAuth Redirect URL. This doesn't matter. We aren't redirecting anywhere. I chose https://discordapp.com
- Give it a category of Chat Bot
- Copy your Client ID

## TWITCH_WATCH_USERNAME
- The username of the twitch user you want to watch

# Adding the Bot to the Discord Server

## Find your Discord Client ID
- Go to https://discordapp.com/developers/applications/ and find your bot
- On the General Information page, you should see Client ID. Copy this
- Go to https://discordapp.com/oauth2/authorize?client_id={ClientID}&scope=bot
- Replace {ClientID} with your Discord client ID
- Note that you must have Manage Server permissions to add the bot to the server

- This is where you would get the notification channel ID