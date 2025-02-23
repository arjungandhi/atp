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

Keeping this simple we will build off the excellent todotxt format. So to keep this simple we will write tools to extend the base todo.sh functionality rather than doing something custom.  

Heres some cool first pass ideas 

1. Ingtegrate with github issues to sync issues with github

### Phases

#### Phase 1: Prototype / Scoping

completed on:

In phase 1 we'll develop the following as a prototype

1. Bare bones task tracking tool
2. Project List
4. Updated Doc for ATP and Project Tracking

Success is all these tools exists and my projects are in that format

#### Phase 2: Beta

Iterate on Phase 1, make changes get to a point where we are happy for 2-3 weeks

Some po

#### Phase 3: Gamma

Code Clean up bug fix make pretty

#### Phase 4: Release

Make a Blog Post on arjungandhi.com, going over project mentality + new tools

