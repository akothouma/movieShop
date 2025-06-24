# 🎬 CineGrid - Your Personal Cinema

CineGrid is a modern web application that allows users to discover, search, and organize movies. Built with Go backend and vanilla JavaScript frontend, it provides a seamless movie browsing experience.

![CineGrid Screenshot](screenshot.png)

## ✨ Features

- **Movie Discovery**: Browse popular and trending movies
- **Search Functionality**: Find movies by title, actor, or keyword
- **Personal Watchlist**: Save movies to watch later
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Movie Details**: View comprehensive information about each movie

## 🛠️ Tech Stack

### Backend
- Go (Golang) 1.24.3
- Standard library HTTP server
- TMDB API integration
- OMDB API integration for additional ratings

### Frontend
- Vanilla JavaScript (ES6+)
- HTML5 & CSS3
- Responsive design with CSS Grid and Flexbox
- Local storage for watchlist persistence

## 🚀 Getting Started

### Prerequisites
- Go 1.24.3 or higher
- TMDB API key (get one at [https://www.themoviedb.org/settings/api](https://www.themoviedb.org/settings/api))

### Installation

1. Clone the repository
   ```bash
   git clone https://github.com/yourusername/cinegrid.git
   cd cinegrid
   ```

2. Create a `.env` file in the root directory with your API keys
   ```
   TMDB_API_KEY=your_tmdb_api_key_here
   ```

3. Install dependencies
   ```bash
   go mod download
   ```

4. Run the backend server
   ```bash
   go run cmd/main.go
   ```

5. In a separate terminal, run the frontend server
   ```bash
   go run cmd/serve_frontend.go
   ```

6. Open your browser and navigate to:
   - Frontend: [http://localhost:3000](http://localhost:3000)
   - Backend API: [http://localhost:8000](http://localhost:8000)

## 📁 Project Structure

```
cinegrid/
├── backend/
│   ├── Handlers/           # HTTP request handlers
│   ├── internals/
│   │   └── models/         # Data models
│   └── utils/              # Utility functions and API clients
├── cmd/
│   ├── main.go             # Backend server entry point
│   └── serve_frontend.go   # Frontend server entry point
├── frontend/
│   └── static/
│       ├── index.html      # Main HTML page
│       ├── index.css       # Styles
│       ├── script.js       # Frontend JavaScript
│       └── placeholder.svg # Placeholder image
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── .env                    # Environment variables (not in repo)
```

## 🔍 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Get a list of movies |
| `/movies` | GET | Get movies with pagination |
| `/movie` | GET | Get detailed information about a specific movie |
| `/search` | GET | Search for movies by query |
| `/debug-search` | GET | Debug endpoint for search functionality |

### Query Parameters

- **page**: Page number for pagination (default: 1)
- **id**: Movie ID for the `/movie` endpoint
- **query**: Search term for the `/search` endpoint

## 💻 Usage

### Browsing Movies
The home page displays a grid of popular movies. Scroll down to load more movies automatically.

### Searching Movies
Use the search bar at the top of the page to find specific movies. Results update as you type.

### Movie Details
Click on any movie card to view detailed information including:
- Title and release year
- Overview/description
- Rating
- Option to add to watchlist

### Watchlist
- Click "Add to Watchlist" on any movie to save it
- Access your watchlist by clicking the "My Watchlist" button
- Mark movies as watched or remove them from your list
- Your watchlist is saved locally and persists between sessions

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [TMDB API](https://developers.themoviedb.org/3) for movie data
- [OMDB API](https://www.omdbapi.com/) for additional ratings
- [Go](https://golang.org/) for the backend language
- [godotenv](https://github.com/joho/godotenv) for environment variable management

---

<div align="center">
  <p>Made with ❤️ by Lorna Akoth</p>
</div>