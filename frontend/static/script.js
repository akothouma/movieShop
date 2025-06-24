// API Configuration
const API_BASE_URL = 'http://localhost:8000'
const TMDB_IMAGE_BASE_URL = 'https://image.tmdb.org/t/p/w500'

// Global variables
let movies = []

let watchlist = JSON.parse(localStorage.getItem("watchlist")) || []
let currentMovie = null

// Fetch movies from API
async function fetchMovies() {
  try {
    const response = await fetch(`${API_BASE_URL}/`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()

    // Transform API data to match frontend format
    movies = data.map(movie => ({
      id: movie.id,
      title: movie.title,
      date: movie.release_date ? movie.release_date.split('-')[0] : 'Unknown',
      poster: movie.poster_path ? `${TMDB_IMAGE_BASE_URL}${movie.poster_path}` : '/placeholder.svg?height=400&width=280',
      description: movie.overview || 'No description available.',
      vote_average: movie.vote_average || 0
    }))

    return movies
  } catch (error) {
    console.error('Error fetching movies:', error)
    // Fallback to empty array if API fails
    movies = []
    showError('Failed to load movies. Please try again later.')
    return movies
  }
}

// Show error message to user
function showError(message) {
  const grid = document.getElementById("moviesGrid")
  grid.innerHTML = `
    <div style="grid-column: 1 / -1; text-align: center; padding: 2rem;">
      <h3 style="color: #e74c3c; margin-bottom: 1rem;">⚠️ Error</h3>
      <p style="color: #666;">${message}</p>
      <button onclick="init()" style="margin-top: 1rem; padding: 0.5rem 1rem; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer;">
        Retry
      </button>
    </div>
  `
}

// Initialize the app
async function init() {
  // Show loading state
  const grid = document.getElementById("moviesGrid")
  grid.innerHTML = `
    <div style="grid-column: 1 / -1; text-align: center; padding: 2rem;">
      <h3 style="color: #666;">🎬 Loading movies...</h3>
    </div>
  `

  // Fetch movies from API
  await fetchMovies()

  // Render the movies
  renderMovies()
  updateWatchlistCount()
}

// Render movies grid
function renderMovies() {
  const grid = document.getElementById("moviesGrid")

  if (movies.length === 0) {
    grid.innerHTML = `
      <div style="grid-column: 1 / -1; text-align: center; padding: 2rem;">
        <h3 style="color: #666; margin-bottom: 1rem;">🎬 No movies available</h3>
        <p style="color: #999;">Check back later for new releases!</p>
      </div>
    `
    return
  }

  grid.innerHTML = movies
    .map(
      (movie) => `
        <div class="movie-card" onclick="openMovieDetails(${movie.id})">
            <img class="movie-poster" src="${movie.poster}" alt="${movie.title}"
                 onerror="this.src='/placeholder.svg?height=400&width=280'">
            <div class="movie-info">
                <h3 class="movie-title">${movie.title}</h3>
                <p class="movie-date">${movie.date}</p>
                ${movie.vote_average ? `<p class="movie-rating">⭐ ${movie.vote_average.toFixed(1)}</p>` : ''}
            </div>
        </div>
    `,
    )
    .join("")
}

// Open movie details modal
function openMovieDetails(movieId) {
  currentMovie = movies.find((m) => m.id === movieId)
  if (!currentMovie) return

  document.getElementById("modalPoster").src = currentMovie.poster
  document.getElementById("modalPoster").onerror = function() {
    this.src = '/placeholder.svg?height=400&width=280'
  }
  document.getElementById("modalTitle").textContent = currentMovie.title
  document.getElementById("modalDate").textContent = `Released: ${currentMovie.date}${currentMovie.vote_average ? ` • ⭐ ${currentMovie.vote_average.toFixed(1)}` : ''}`
  document.getElementById("modalDescription").textContent = currentMovie.description

  const addBtn = document.getElementById("addToWatchlistBtn")
  const isInWatchlist = watchlist.some((item) => item.id === currentMovie.id)

  if (isInWatchlist) {
    addBtn.textContent = "Already in Watchlist"
    addBtn.disabled = true
  } else {
    addBtn.textContent = "Add to Watchlist"
    addBtn.disabled = false
  }

  document.getElementById("movieModal").style.display = "block"
}

// Close movie details modal
function closeModal() {
  document.getElementById("movieModal").style.display = "none"
  currentMovie = null
}

// Add movie to watchlist
function addToWatchlist() {
  if (!currentMovie) return

  const isAlreadyInWatchlist = watchlist.some((item) => item.id === currentMovie.id)
  if (isAlreadyInWatchlist) return

  watchlist.push({
    ...currentMovie,
    watched: false,
    addedDate: new Date().toISOString(),
  })

  localStorage.setItem("watchlist", JSON.stringify(watchlist))
  updateWatchlistCount()

  const addBtn = document.getElementById("addToWatchlistBtn")
  addBtn.textContent = "Added to Watchlist!"
  addBtn.disabled = true

  setTimeout(() => {
    closeModal()
  }, 1000)
}

// Open watchlist modal
function openWatchlist() {
  renderWatchlist()
  document.getElementById("watchlistModal").style.display = "block"
}

// Close watchlist modal
function closeWatchlist() {
  document.getElementById("watchlistModal").style.display = "none"
}

// Render watchlist items
function renderWatchlist() {
  const container = document.getElementById("watchlistItems")

  if (watchlist.length === 0) {
    container.innerHTML = `
            <div class="empty-watchlist">
                <h3>Your watchlist is empty</h3>
                <p>Add some movies to get started!</p>
            </div>
        `
    return
  }

  container.innerHTML = watchlist
    .map(
      (movie) => `
        <div class="watchlist-item ${movie.watched ? "watched-item" : ""}">
            <img class="watchlist-poster" src="${movie.poster}" alt="${movie.title}">
            <div class="watchlist-info">
                <div class="watchlist-movie-title">
                    ${movie.title}
                    ${movie.watched ? '<span class="watched-badge">Watched</span>' : ""}
                </div>
                <div class="watchlist-movie-date">${movie.date}</div>
            </div>
            <div class="watchlist-actions">
                ${
                  !movie.watched
                    ? `
                    <button class="action-btn watched-btn" onclick="markAsWatched(${movie.id})">
                        Mark Watched
                    </button>
                `
                    : ""
                }
                <button class="action-btn remove-btn" onclick="removeFromWatchlist(${movie.id})">
                    Remove
                </button>
            </div>
        </div>
    `,
    )
    .join("")
}

// Mark movie as watched
function markAsWatched(movieId) {
  const movieIndex = watchlist.findIndex((item) => item.id === movieId)
  if (movieIndex !== -1) {
    watchlist[movieIndex].watched = true
    localStorage.setItem("watchlist", JSON.stringify(watchlist))
    renderWatchlist()
  }
}

// Remove movie from watchlist
function removeFromWatchlist(movieId) {
  watchlist = watchlist.filter((item) => item.id !== movieId)
  localStorage.setItem("watchlist", JSON.stringify(watchlist))
  updateWatchlistCount()
  renderWatchlist()
}

// Update watchlist count badge
function updateWatchlistCount() {
  document.getElementById("watchlistCount").textContent = watchlist.length
}

// Close modals when clicking outside
window.onclick = (event) => {
  const movieModal = document.getElementById("movieModal")
  const watchlistModal = document.getElementById("watchlistModal")

  if (event.target === movieModal) {
    closeModal()
  }
  if (event.target === watchlistModal) {
    closeWatchlist()
  }
}

// Initialize the app when page loads
init()
