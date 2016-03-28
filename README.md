# SteamGroupMemberStats [![Build Status](https://travis-ci.org/AzureByte/SteamGroupMemberStats.svg?branch=master)](https://travis-ci.org/AzureByte/SteamGroupMemberStats)


An attempt at visualizing Steam Community data. This project aims at helping moderators and group owners better know the games their members are playing and how active their community is.

Early stage example :
![enter image description here](http://i.imgur.com/LiQttDu.jpg)



*This project was started during the 2016 Mumbai Hackathon*


## Installation

### Windows
Download the bin folder and run httpserver.exe

### Other OS
To compile from source, you will need [Go](https://golang.org/dl/) installed.

From the goserver folder, run the following command:
>go run httpserver.go

Regardless of the OS, it will create a server on port 8888

Visit [http://localhost:8888/](http://localhost:8888/) to view the App. Then click on any of the communities listed to get a breakdown of the games people are playing.

Note : It currently only shows hardcoded data obtained from the [Steam Stats page](http://store.steampowered.com/stats/).