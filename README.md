# aoe-balance-bot
A Discord bot written in Go to balance AoE2 DE lobbies using in-game id's. 

# Commands
- `!balance`: Create two teams of players with minimal ELO difference in an AoE2DE lobby using its provided game id (`!balance aoe2de://0/xxxxxxxxx`). This command will also run if a lobby id is pasted into a channel without the command prefix.
- `!help`: Display all available commands

# Setup
The bot requires the following environment variables to be set:
- `token`: Discord bot token

And provides the following optional parameters:
- `cacheExpiryHours` How long player data should be retained in the cache (default: `24` hours).
- `cacheMaxSize` How many players data can be stored at one time before stale data will be removed (default `20` entities).
- `port` Port used for heartbeat (If this is not provided no server will be started).
- `trustInsecureCertificates` Disable certificate validation when making API calls to get game data.
- `whitelistedChannels` List of channel id's. If this is set, the bot will only react to messages in these channels (defaults to all channels).