// API Configuration
const API_BASE_URL = 'http://localhost:8000'
const TMDB_IMAGE_BASE_URL = 'https://image.tmdb.org/t/p/w500'

// Global variables
let movies = []
let currentPage = 1
let totalPages = 1
let isLoading = false

let watchlist = JSON.parse(localStorage.getItem("watchlist")) || []
let currentMovie = null

// Fetch movies from API with pagination
async function fetchMoviesPage(page) {
  try {
    console.log(`Fetching movies page ${page}...`)
    isLoading = true
    
    const response = await fetch(`${API_BASE_URL}/movies?page=${page}`)
    if (!response.ok) {
      const errorText = await response.text()
      console.error(`HTTP error! status: ${response.status}, message: ${errorText}`)
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const data = await response.json()
    console.log('API response:', data)
    
    // Update pagination info
    currentPage = data.page
    totalPages = data.total_pages
    
    // Transform API data to match frontend format
    const newMovies = data.results.map(movie => ({
      id: movie.id,
      title: movie.title,
      date: movie.release_date ? movie.release_date.split('-')[0] : 'Unknown',
      poster: movie.poster_path ? `${TMDB_IMAGE_BASE_URL}${movie.poster_path}` : '/placeholder.svg?height=400&width=280',
      description: movie.overview || 'No description available.',
      vote_average: movie.vote_average || 0
    }))
    
    isLoading = false
    return {
      movies: newMovies,
      currentPage,
      totalPages
    }
  } catch (error) {
    console.error('Error fetching movies:', error)
    isLoading = false
    showError('Failed to load movies. Please try again later.')
    return { movies: [], currentPage: 1, totalPages: 1 }
  }
}

// Load more movies (next page)
async function loadMoreMovies() {
  if (isLoading || currentPage >= totalPages) return
  
  isLoading = true
  const loadMoreBtn = document.getElementById('loadMoreBtn')
  
  if (loadMoreBtn) {
    loadMoreBtn.innerHTML = `
      <span class="btn-content">
        <span>Loading Movies</span>
        <span class="loading-dots">
          <span class="loading-dot"></span>
          <span class="loading-dot"></span>
          <span class="loading-dot"></span>
        </span>
      </span>
    `
    loadMoreBtn.disabled = true
  }
  
  try {
    const result = await fetchMoviesPage(currentPage + 1)
    movies = [...movies, ...result.movies]
    
    // Add a slight delay for better UX
    setTimeout(() => {
      renderMovies(false) // Append mode
      updatePaginationUI()
      isLoading = false
      
      // Smooth scroll to show new content
      const lastOldMovieIndex = (currentPage - 1) * 20
      const firstNewMovie = document.querySelectorAll('.movie-card')[lastOldMovieIndex]
      if (firstNewMovie) {
        window.scrollTo({
          top: firstNewMovie.offsetTop - 120,
          behavior: 'smooth'
        })
      }
    }, 800)
  } catch (error) {
    console.error('Error loading more movies:', error)
    isLoading = false
    
    if (loadMoreBtn) {
      loadMoreBtn.innerHTML = `
        <span class="btn-content">
          <span>Error Loading Movies</span>
          <span class="btn-icon">⚠️</span>
        </span>
      `
      
      setTimeout(() => {
        updatePaginationUI()
      }, 2000)
    }
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

  // Fetch first page of movies
  const result = await fetchMoviesPage(1)
  movies = result.movies
  
  // Render the movies
  renderMovies(true) // Replace mode
  updatePaginationUI()
  updateWatchlistCount()
}

// Render movies grid
function renderMovies(replace = true) {
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

  const moviesHTML = movies
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
  
  if (replace) {
    grid.innerHTML = moviesHTML
  } else {
    // Append new movies to existing grid
    const tempDiv = document.createElement('div')
    tempDiv.innerHTML = moviesHTML
    while (tempDiv.firstChild) {
      grid.appendChild(tempDiv.firstChild)
    }
  }
  
  // Add pagination controls after the grid
  if (replace) {
    const paginationContainer = document.createElement('div')
    paginationContainer.id = 'paginationContainer'
    paginationContainer.className = 'pagination-container'
    grid.parentNode.insertBefore(paginationContainer, grid.nextSibling)
  }
}

// Update pagination UI with enhanced load more button
function updatePaginationUI() {
  const container = document.getElementById('paginationContainer')
  if (!container) return
  
  const percentLoaded = Math.min(100, (currentPage / totalPages) * 100).toFixed(0)
  
  container.innerHTML = `
    <div class="pagination-progress">
      <div class="pagination-progress-bar" style="width: ${percentLoaded}%"></div>
    </div>
    <div class="pagination-info">
      Showing ${movies.length} of approximately ${totalPages * 20} movies (${percentLoaded}% loaded)
    </div>
    <button id="loadMoreBtn" class="load-more-btn" ${currentPage >= totalPages ? 'disabled' : ''}>
      <span class="btn-content">
        ${currentPage >= totalPages ? 
          '<span>All Movies Loaded</span> <span class="btn-icon">✓</span>' : 
          '<span>Discover More Movies</span> <span class="btn-icon">↓</span>'}
      </span>
    </button>
  `
  
  const loadMoreBtn = document.getElementById('loadMoreBtn')
  if (loadMoreBtn && currentPage < totalPages) {
    loadMoreBtn.addEventListener('click', loadMoreMovies)
    
    // Add hover animation for the icon
    loadMoreBtn.addEventListener('mouseenter', () => {
      const icon = loadMoreBtn.querySelector('.btn-icon')
      if (icon) {
        icon.style.animation = 'bounce 0.5s ease infinite'
      }
    })
    
    loadMoreBtn.addEventListener('mouseleave', () => {
      const icon = loadMoreBtn.querySelector('.btn-icon')
      if (icon) {
        icon.style.animation = ''
      }
    })
  }
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
