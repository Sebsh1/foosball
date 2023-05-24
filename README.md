### A service for handling you office or club needs for logging played matches
Does your office have a foosball table or table tennis?
Then you and your colleagues have no doubt discussed coming up with a rating system for your favourite game to see they stack up against each other. \
This service does that for you, giving you more time to do actual work. \
**Beware:** You need a healthy work environment, or else ranking your you and your colleagues on any measure might lead to people feeling inadeqaute or unwelcome.

### Features
This repository has a complete backend and REST API for logging matches to a database including
- Organizing users into organizations
- Calculating ratings using varius methods (Elo, Weighted & Glicko2)
- Keeping track of player statistics
- Leaderboards
- Self-host or use the API at matchlog.com
- 
All you have to do is deploy it, play your favorite games and log the results. 

**Notice:** Your employees might instead start thinking about creating a frontend for this service to throw onto a screen near the games they play. \
I currently have a frontend in its infancy as I try to learn how to frontend. :)

### Setup
- Define a <code>.env</code> file, see <code>.env.example</code> for required variables.
- Run <code>go run main.go serve</code> to start the service with the environment defined in <code>.env</code>.

### TODO by priority
- TopX lists (rating, win/loss ratio, total matches, etc.)
- New rating methods: Glicko2 & more?
- Organizing and managing tournaments
- Organizing and managing leagues (maybe just a subset of tournaments if supporting round-robin tournaments?) 
- Docker image for running locally
- API Docs
- a frontend :)
