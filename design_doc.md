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


## Project Tracking Tool

Needs work revist after 


## Task Execution Tool 

A task is just a thing I need to do. It can be anything related to projects or not. 

Tasks have the following key info

1. Name
2. Due Date
4. Expected Duration
5. Actual Duration
6. Actual Completion Date
7. Dependencies List[Task] - Future
8. Labels (key:values to extend the task with other systems)

### Existing Tools I researched

1. Open Project
2. Plane
3. Click up 
4. Github Issues Based
5. Todoist 
6. Task Warrior
7. todo.txt

### Task Tracking Design 

1. list all tasks 
2. start task
3. stop a task 
4. manual priority order
5. project label


we'll build this app off of the [todo.txt](https://github.com/1set/todotxt) format, its a simple text based format with lots of community support. It also makes things like future integration much easier.

#### File format 

https://github.com/todotxt/todo.txt


### Access

This tool needs both ios and computer access in the long term but maybe for now we just start off with laptop access 


### Phases

#### Phase 1: Prototype / Scoping

completed on:

In phase 1 we'll develop the following as a prototype

1. Bare bones task tracking tool
2. Project List
3. Template For Project Docs
4. Updated Doc for ATP and Project Tracking

Success is all these tools exists and my projects are in that format

#### Phase 2: Beta

Iterate on Phase 1, make changes get to a point where we are happy for 2-3 weeks

##### Beta Ideas
1. Alarm once we pass estimated time 
2. Hard Expiry on tasks 
3. Github Issue Integration
4. Automation to add tasks on a regular basis
5. Journal? extract todos from journal into app

#### Phase 3: Gamma

Code Clean up bug fix make pretty

#### Phase 4: Release

Make a Blog Post on arjungandhi.com, going over project mentality + new tools

