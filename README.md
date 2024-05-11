# unsocial-network

## Description

This project aims to develop a Facebook-like social network with various features including followers, profiles, posts, groups, notifications, chats, and more. Below are the key requirements and instructions for implementation.

## Implementation
Backend: Go,<br>
Frontend: JavaScript & Vue.js<br>
Database: SQLite

## Run the project

Docker:
Use docker to run the project
```bash
docker-compose up --build
```

and navigate to http://localhost:8080 or click on the link in the terminal.

OR run the project manually:

there is also a small script for your convenience:
```bash
sh start-script.sh
```

## Documentation

![Home page](/screenshots/unsocial-network_home.png "Home page")

### Authentication

Registration and login forms with sessions and cookies.
Required registration information includes email, password, first name, last name, and date of birth.
Optional fields include avatar, nickname and about me.

### Followers

Users can follow and unfollow other users while navigating the application. Implementation for follow request functionality.

### Profile

Profiles contain user information, activity and and lists for follower & following. Profiles can be public or private, with an option to toggle privacy settings. Public profiles don't require a follow request.

![Profile page](/screenshots/unsocial-network_profile.png "Profile page")

### Posts

Users can create posts and comment on other user's posts.
Post privacy settings include public, private and almost private (friends).

### Groups

Users can create groups with titles and descriptions.
Groups support invitations and requests for memberships.
Group members can create events within the group.

![Groups page](/screenshots/unsocial-network_groups.png "Groups page")

### Chats

Private messaging between users. Groups have a common chat room for members.

### Notifications

Users receive notifications for follow requests, group invitations, group membership requests, chat messages and group events.

## Contributors
[4stroPhysik3r](https://github.com/4stroPhysik3r)<br>
Freyby<br>
KristjanM