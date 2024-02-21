# Tinycare-tui

Small terminal app that shows git commits from the last 24 hours and week, current weather, some self care advice, and you current todo list tasks
inspired by @notwaldorf [tiny-care-terminal](https://github.com/notwaldorf/tiny-care-terminal)

## Installation
```
go install github.com/DMcP89/tinycare-tui@latest
```

## TO-DOs
- [x] Allow for focusing on each box
- [x] Expand on self care reminders
- [x] Remove twitter scraping code
- [x] Replace hardcoded values with env variables
- [x] Refresh the view on 'r'
- [x] Have views refresh in 30 second intervals
- [x] Have 'q' exit the app
- [x] Optimize with go routines
- [x] Allow for multiple git repo locations
- [x] Allow for local todo list for todays tasks
- [x] Have task view show completed tasks as well
- [x] Provide option to pull commits from github instead of from local repos
- [ ] Add error handling for missing environment variables
- [ ] Performance tuning
- [ ] Logging
- [ ] Refactoring
- [ ] Convert time on commits to days when >24 hours
- [ ] Refactor GitHub interactions to use go-hithub
- [x] Write installation guide


## About
I started this project to accomplish a few different goals
1. Teach myself the basics of the Go language
2. Practice leveraging generative AI tools like copilot and chatgpt for development
3. Create a fun terminal app
