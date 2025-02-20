# GATOR - Blog Aggregator
Gator will download posts from your favourite RSS feeds and you can read them right in your terminal.
## Requirements
Gator requires that you've installed the Go programming language and have a postgres server available. 
This config will later include the current username, however this will be added by Gator.
## Commands and Use
Gator uses the following commands:
`register [username]` Register your user - argument mandatory
`login [username]` Login your user - argument mandatory
`reset` Reset user database - use with caution
`users` Get a list of users
`agg [update frequency]` Start fetching posts from feeds - argument mandatory ex:(2s, 5m, 1hr)
`addfeed [feedname] [feed url]` Add a feed - arguments both mandatory
`feeds` List available feeds
`follow [feed url]` Follow a feed by URL - argument mandatory
`following` List your followed feeds
`unfollow [feed url]` Unfollow a feed by URL - argument mandatory
`browse [post count]` Display a count of posts - argumetn optional, default is 2
## Installation
Installation is simple. Install via go install
`$go install github.com/MattInReality/gator`
Or download source and run go build to create the binary yourself.
Gator is configured by having a .gatorconfig.json file in your home directory. You will need to add your postgres connection url to the file as shown here:
```json
.gatorconfig.json
{
    "db_url": "postgres://yourname:yourpassword@dbIP:db:PORT/gator"
}
    
```

