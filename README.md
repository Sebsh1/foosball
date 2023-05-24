### Solving Office & Club Needs for Logging Played Matches
Does your office have a foosball table or table tennis?
Then you and your colleagues have no doubt discussed coming up with a rating system for your favourite game to see they stack up against each other. \
This service does that for you, giving you more time to do actual work. \
**Beware:** You need a healthy work environment, or else ranking your you and your colleagues on any measure might lead to people feeling inadeqaute or unwelcome.

### Features
This repository has a complete backend and REST API for logging matches to a database including
- Organizing users into organizations
- Calculating ratings using varius methods (Elo, Weighted & Glicko2)
- Keeping track of player statistics

All you have to do is deploy it, play your favorite games and log the results. 

**Notice:** Your employees might instead start thinking about creating a frontend for this service to throw onto a screen near the games they play. \
I currently have a frontend in its infancy as I try to learn how to frontend. :)

### Setup
- Define a <code>.env</code> file, see <code>.env.example</code> for required variables.
- Run <code>go run main.go serve</code> to start the service with the environment defined in <code>.env</code>.

### Contributing
If you would like to contribute with features/fixes/etc have a look at the open [Issues](https://github.com/Sebsh1/matchlog/issues) for my envisioned future features or create your own. I am also open for PRs.
