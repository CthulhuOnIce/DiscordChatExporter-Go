# Usage
discorddumper \[user_id\] \[--bot?\] --channels=0,1,2 --guilds=0,1,2 --media=

# BluePrint
- Client
    - Media[]
    - CachedGuilds[]
    - GetGuild()
    - EnumGuilds()

- Guild
    - CachedChannels\[Channel\]
    - GetChannel()
    - EnumChannels()
    - client Client

- Channel
    - bool IsVoice
    - guild Guild