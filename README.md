
<p  align="center">

<img  src="https://user-images.githubusercontent.com/30529572/72455010-fb38d400-37e7-11ea-9c1e-8cdeb5f5906e.png"  />

<h2  align="center"> Youtube Integration API </h2>

<h4  align="center"> The API can be used to integrate aToutube Channel's Playlist and Videos on their websites. <h4>

</p>

---

[![DOCS](https://img.shields.io/badge/Documentation-see%20docs-green?style=flat-square&logo=appveyor)](https://documenter.getpostman.com/view/10749950/SzfAzS8f?version=latest)

  

## Functionalities

- [x]  Endpint to fetch playlist and videos not in ay playlist

- [x] Endpoint to fetch videos of a particular playlist 

- [x] Endpoint to update the database


  

<br>

  
  

## Instructions to run

  

* Pre-requisites:

- Golang 1.14
- MongoDB Go Driver 1.3
- Gin

*  Directions to build

```bash

go build

```

  

* Varibles required 

  

```bash

API_KEY="YOUR API KEY"
CHANNEL_ID="CHANNEL ID"
USERNAME="DATABASE ADMIN USERNAME"
PASSWORD="PASSWORD FOR THE DATABSE"

```
* Database Stucture
```
Database:main -> collections: Playlists, Videos
Database:Logs -> collections: UPDATED
```
  

<br>

  

## Contributors

  

*  [Mayank Kumar](https://github.com/mayankkumar2)


  
  
  

<br>

<br>

  

<p  align="center">

Made with :heart: by DSC VIT

</p>