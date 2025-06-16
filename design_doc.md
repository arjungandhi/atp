# ATP Design Doc

# Overview
Task Planning and Scope Creep is a pretty difficult challenge for me, this doc aims to define a system + tools to help make that easier. 


# Background / Motivations

My life is full of projects, I consider most things to be a project (work, relationships, dogs, family, exercise) all of it falls in the project category. 

I also SUCK at finishing projects. 

I would like to get better at this. 

I recently was in a talk for pattern from [Sean Behr](https://www.linkedin.com/in/seanbehr) CEO of [Fountian](fountain.com). Where he talked about turning fountian into a "high velocity product factory"
very buzz wordy I know. This is the basis for the system described below

# System Overview

This doc is really information about 2 projects in one but for design purposes its simpler to group them as one. 

## Projects (Planning Scope Release)

Projects are the current highest level item in the Arjun Planning World. 

### Phases of Projects
1. Scoping & Design and Research & Prototypes- What the hell is thing, what do people / or me need and how do we get them to it.
2. Beta (Build the damn thing and Use) - In this phase its all about building the minimal thing that achieves the goal if you're smart you'll have the "customers" working side by side with you to help develop the right thing
3. GA (Release it) - In this phase you "scale up" and "release the product" to the wild. You spend some time ironing out bugs and acquiring more customers. And then its gone. No adding "just one feature" to increase the scope.
4. Finish - In this phase my goal is to either 1. Hand this off to some one who can run / maintain it better than I can or 2. Call it done. Its hard to keep maintains up on these things (naturally not everything falls in this category but I feel like a bulk of my ideas can be completed). Would love to start writing these changes up in a blog or something at some point in the future


### Rules for Projects
1. I may have N (probably 3) (Not including the non negotiable's relationships, eating, exercise etc) active projects at any one time. Any more than that is probably too much.
2. I will maintain a project list which to write down future projects and keep track of which projects are over the line
3. Finish it. If I choose to go to phase 2 I must make it to phase 4. No 20% at the end for me
4. Killing projects is fine and encouraged! But you may only do it in between phase 1 & 2. 
5. Have fun mf!

### Measuring Success
If this is successful 2025 we should do alot more stuff. If this system is useful it'll become second nature. I will check in at each quarter to see if I actually am 

1. Using the system. 
2. Completing Projects. 
3. Having Fun.

To measure this I need to keep track of projects some how. 


## Project Tracking Design


I can do the following actions with projects

1. Add a new project idea - create a new entry in projects.txt
2. Make a project active - Prioritize the project in projects.txt, set its phase, if either select a repo or attach it to an existing repo (if needed), clone that repo locally (if needed)
3. Edit a project doc - open the doc.txt in the project directory
4. Deactivate a project - remove the repo from my computer + deprioritize the entry in projects.txt
5. Finish a project - remove the repo from my computer + move the project entry from projects.txt -> finished_projects.txt
6. Delete a project idea - delete a entry in projects.txt
7. Set the phase of a project - change the phase label in a project 

The project and finished_project files will follow the [todo.txt](https://github.com/1set/todotxt) format. Its simple and expandable and parseable by a simple text editor as well as has good connection with other tools  

## Task Execution Tool

Keeping this simple we will build off the excellent todotxt format. So to keep this simple we will write tools to extend the base todo.sh functionality rather than doing something custom. For now I will just use the todo command and will only add functionality to edit the todo.txt easier

### Features

#### Recurring Tasks 
Reocurring tasks will be stored in a recur.txt file 
Recurring tasks will be similar in format to a crontab file. 
It will have the following format

Example:
```
@daily do task @tag +context
@weekly do task @tag +context
0 0 * * * do task @tag +context
```

We can then run `atp recur` to look at the tasks that are in the recur.txt file and add them to the todo.txt file if needed
We will need to be smart about this and not add the task for a period multiple times.

##### Implementation Details
- Recurring tasks support both simple formats (`@daily`, `@weekly`, `@monthly`) and full cron format (`0 9 * * 1`)
- Generated todos include a `recur:YYYY-MM-DD` label to track when they were generated and prevent duplicates
- The system checks existing todos before generating new ones to avoid creating duplicate tasks for the same date
- `atp recur edit` opens the recur.txt file for editing recurring task templates
- The recurring task system integrates with the existing todo.txt format and file structure

#### Reminder Tasks
Reminder tasks are one-time future tasks that should only appear in the active todo list on or after their specified reminder date. This is useful for tasks that need to be done in the future but aren't relevant until that date arrives.

Examples:
- `evaluate if system is still working cancel if not remind:2025-07-15`
- `cancel renters insurance remind:2025-11-01 @home +personal`

##### Storage Format
- Reminder tasks are stored in `$ATP_DIR/todo/reminders.txt` using standard todo.txt format
- Each task includes a `remind:YYYY-MM-DD` label specifying when it should become active
- Tasks can include all standard todo.txt elements: priority, projects (+), contexts (@), and other labels

##### File Structure
- `$ATP_DIR/todo/reminders.txt` - Future reminder tasks waiting to be activated
- Tasks remain in `reminders.txt` until their reminder date arrives
- On the reminder date, tasks are moved from `reminders.txt` to `todo.txt` and removed from reminders

##### CLI Commands
- `atp todo remind add [task]` - Add a new reminder task with interactive date prompt
- `atp todo remind edit` - Edit the reminders.txt file directly
- `atp todo remind list` - List all pending reminder tasks sorted by date
- `atp todo remind process` - Process reminders for today (move due reminders to active todos)
- `atp todo remind process [date]` - Process reminders for a specific date

##### Implementation Details
- Reminder tasks use the existing Todo struct and parsing logic
- The `remind:YYYY-MM-DD` label is used to determine activation date
- Tasks are moved to active todos and deleted from reminders.txt in a single operation
- No duplicate checking needed since tasks are removed from reminders once processed
- Integration with existing file backup and error handling patterns
- Compatible with all existing todo.txt tooling and formats

## Architecture Changes

### Package Structure Refactoring
The codebase has been refactored from a single `pkg` package to a more modular structure:

- `cmd/atp/` - Main CLI application with subcommands organized into separate files
- `cmd/atp/cli/` - CLI command implementations and utilities  
- `project/` - Project management functionality (moved from `pkg`)
- `todo/` - Todo.txt format parsing and manipulation (moved from `pkg/todo`)
- `repo/` - Repository management and discovery (moved from `pkg/repo`)

### File Organization Improvements
- Commands are now organized into separate files (`todo.go`, `project.go`, `load.go`)
- Utility functions consolidated into `load.go` for better organization
- Test files moved alongside their corresponding implementation files
- Improved separation of concerns with dedicated directories for each major component
