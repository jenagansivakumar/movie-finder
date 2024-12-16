# Movie Recommender
 
## Overview

The Movie Recommender app was created to address challenges such as web scraping and enhanced response times. Using The Movie Database (TMDb) API, this app fetches popular movie data, implements rate limiting, and utilises caching with Redis to speed up responses.

The app is containerised using Docker, including a Redis container, making it portable and easy to deploy.
## Features

- **Rate Limiting**: Uses `golang.org/x/time/rate` and Redis to prevent excessive requests to the TMDb API (1 request per 5 seconds).
- **Caching**: Data from TMDb is cached in Redis to avoid repeated API calls and improve response times.
- **Containerisation**: The entire system, including the Movie Recommender app and Redis, is containerised using Docker.

## Technologies Used

- **Go**: The backend language for the application.
- **Redis**: For caching API responses and rate limiting.
- **Docker**: To containerise both the application and Redis for easy deployment.
- **Rate Limiting**: Implemented using `golang.org/x/time/rate` along with Redis for storing user request limits.

## How It Works

1. The app fetches popular movie data from TMDb.
2. It applies a rate limiter to avoid overwhelming the API (1 request every 5 seconds).
3. The data is cached in Redis for quicker access on subsequent requests.
4. Docker is used to containerise the entire system, making it easy to deploy and run anywhere.

## Installation & Setup

### 1. Clone the repository:

```bash
git clone https://github.com/yourusername/movie-recommender.git
cd movie-recommender
```

### 2. Create the .env file:
```
API_KEY=your_api_key_here
```

### 3. Run with Docker Compose:
Make sure Docker and Docker Compose are installed. Then, run:
```
docker compose up --build
```
This comment will build the Docker images and start both the Movie Recommender app and Redis containers.

### 4. Access the app:
Once the container are up and running, visit http://localhost:8080 in your browser. Refresh the page a few times to see the logs and check the caching behaviour.

## How to Test
1. Once the app is running, check the logs:
- "cache miss"
- "fetching data from API"
- "setting cached data"
- "using cache"
2. Logs wwill show when the data is fetched from the API and when it's served from the cache.
