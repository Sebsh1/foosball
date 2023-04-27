### A service for handling you office or club needs for logging played matches
Does your office have a foosball table or table tennis?
Then they have no doubt discussed coming up with a rating system for their matches to see how they stack up against each other. \
This service does that for you, giving you more time to do actual work. 

### Features
This repository has a complete backend and REST API for logging matches to a database and calculating rating, leaderboards and player stats. \
All you have to do is deploy it, play your favorite games and log the results. 

**Notice:** Your employees might instead start thinking about creating a frontend for this service to throw onto a screen near the table.
I currently have a frontend in its infancy as I try to learn react. :)

### Setup
Run <code>go run main.go serve</code> to start the service with the config defined in <code>config.yaml</code>.

### TODO
- Authguard on endpoints
- Better stats collection 
- TopX lists of various stats
- Organizing and managing tournaments  
- API Docs
- a frontend :)
