
# movie-finder-api

Movie Finder API

A Go-based API for fetching movie recommendations with rate limiting, integrating The Movie Database (TMDb) API for real-time movie data. Fully Dockerised for easy deployment and testing.



## Features

- Movie Recommendations: Fetch movies dynamically using query parameters (e.g., `/recommendations?genre=Sci-Fi`).
- Rate Limiting: IP-based rate limiter (1 request every 5 seconds) to control traffic and mitigate scraping.
- TMDb API Integration: Retrieves movie data such as title, genre, and ratings.
- Health Check: `/health` endpoint for monitoring server status.
- Dockerised Deployment: Simplified setup with Docker.

---
## Requirements
- Go (version 1.20+)
- Docker
- TMDb API Key

---
## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/jenagansivakumar/movie-finder.git
cd movie-finder
```

### 2. Add TMDb API Key
- Create a `.env` file in the project root:
  ```
  TMDB_API_KEY=your_api_key_here
  ```

### 3. Build and Run with Docker
```bash
docker build -t movie-recommender .
docker run -p 8080:8080 --env-file .env movie-recommender
```

---
## Endpoints

### 1. `/health`
- Method: `GET`
- Response: `OK`

### 2. `/recommendations?genre=<genre>`
- Method: `GET`
- Query Parameter:
  - `genre`: Movie genre or keyword to search.
- Response: JSON array of movies.

---
## Demo
- https://drive.google.com/file/d/1BLfOVsc4-rdOqWF9G-8i_NFZ8Km0UgCS/view?usp=drive_link <----- Link to video demo

---
