### Solving Office & Club Needs for Logging Played Matches
Does your office have a foosball table, table tennis or a similar game they enhoy playing?
Then you and your colleagues have no doubt discussed coming up with a rating system for your favourite game to see they stack up against each other. \
This service does that for you, giving you more time to do actual work. \
**Beware:** You need a healthy work environment, or else ranking you and your colleagues on any easily viewable measure might lead to people feeling inadeqaute or unwelcome.

### Features
This repository has a complete backend and REST API for logging matches to a database including
- Organizing users into organizations
- Calculating ratings using varius methods (Elo, Weighted & ~~Glicko2~~)
- Keeping track of player statistics including various leaderboards

**Notice:** Your employees might instead start thinking about creating a frontend for this service to throw onto a screen near the games they play. \
I currently have a frontend in its infancy as I try to learn how to frontend. :)

### Contributing
If you would like to contribute with features/fixes/etc have a look at the open [Issues](https://github.com/Sebsh1/matchlog/issues) for my envisioned future features or create your own. I am also open for PRs.

### API
The API exposed by the service is documented at https://sebsh1.github.io/matchlog/ and is intended for use by some frontend application, since all endpoints require a valid JWT issued by the /login endpoint, which periodically expires. \
In the future I might add separate endspoints requiring an API key instead, but not for now.
