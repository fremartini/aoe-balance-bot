# aoe-balance-bot
A Discord bot written in Go to balance AoE2 DE lobbies using in-game id's. 

# Commands
- `!balance`: Create two teams of players with minimal ELO difference in an AoE2DE lobby using its provided game id (`!balance aoe2de://0/xxxxxxxxx`).
- `!help`: Display all available commands

# Setup
The bot requires the following environment variables to be set:
- `token`: Discord bot token
- `logLevel` Level to filter logs displayed in the console 
    - 0: `INFO`
    - 1: `WARN`
    - 2: `FATAL`