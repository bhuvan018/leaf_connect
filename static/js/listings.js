// listings.js - Listing-related functionality

// State
let currentListings = [];
let favorites = [];

/**
 * Fetch and display all listings
 * @param {Object} filters - Optional filters to apply
 */
async function fetchListings(filters = {}) {
  try {
    // Build query string from filters
    const queryParams = new URLSearchParams();
    
    if (filters.userId) queryParams.append('userId', filters.userId);
    if (filters.type) queryParams.append('type', filters.type);
    if (filters.plantType) queryParams.append('plantType', filters.plantType);
    if (filters.location) queryParams.append('location', filters.location);
    
    const url = `/api/listings${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
    
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('Failed to fetch listings');
    }
    
    const listings = await response.json();
    currentListings = listings;
    
    // Fetch user's favorites for badge display
    await fetchFavorites();
    
    // Display listings
    displayListings(listings);
    
    return listings;
  } catch (error) {
    console.error('Error fetching listings:', error);
    displayError('Failed to load listings. Please try again later.');
    return [];
  }
}

/**
 * Display listings in the listing container
 * @param {Array} listings - Array of listing data
 */
function displayListings(listings) {
  const container = document.getElementById('listings-container');
  if (!container) return;
  
  // Clear existing content
  container.innerHTML = '';
  
  if (listings.length === 0) {
    container.innerHTML = '<p class="text-center">No listings found. Try adjusting your filters.</p>';
    return;
  }
  
  // Create grid container
  const grid = document.createElement('div');
  grid.className = 'grid';
  
  // Add listing cards
  listings.forEach(listing => {
    const card = createListingCard(listing);
    grid.appendChild(card);
  });
  
  container.appendChild(grid);
}

/**
 * Handle search form submission
 * @param {Event} event - Form submit event
 */
async function handleSearch(event) {
  event.preventDefault();
  
  const searchInput = document.getElementById('search-input');
  if (!searchInput) return;
  
  const query = searchInput.value.trim();
  if (!query) {
    // If empty search, show all listings
    await fetchListings();
    return;
  }
  
  try {
    const response = await fetch(`/api/listings/search?q=${encodeURIComponent(query)}`);
    if (!response.ok) {
      throw new Error('Search failed');
    }
    
    const results = await response.json();
    currentListings = results;
    
    // Display search results
    displayListings(results);
  } catch (error) {
    console.error('Search error:', error);
    displayError('Search failed. Please try again.');
  }
}

/**
 * Handle filter changes
 */
function handleFilterChange() {
  const typeFilter = document.getElementById('type-filter');
  const plantTypeFilter = document.getElementById('plant-type-filter');
  const locationFilter = document.getElementById('location-filter');
  
  const filters = {};
  
  if (typeFilter && typeFilter.value) {
    filters.type = typeFilter.value;
  }
  
  if (plantTypeFilter && plantTypeFilter.value) {
    filters.plantType = plantTypeFilter.value;
  }
  
  if (locationFilter && locationFilter.value) {
    filters.location = locationFilter.value;
  }
  
  // Fetch listings with filters
  fetchListings(filters);
}

/**
 * Create a new listing
 * @param {Event} event - Form submit event
 */
async function createListing(event) {
  event.preventDefault();
  
  const form = event.target;
  const formData = new FormData(form);
  
  // Convert form data to JSON
  const listingData = {};
  formData.forEach((value, key) => {
    // Handle numeric values
    if (key === 'price') {
      listingData[key] = parseFloat(value);
    } else if (key === 'images') {
      // For demo, use a placeholder image if none provided
      listingData[key] = value ? [value] : [];
    } else {
      listingData[key] = value;
    }
  });
  
  try {
    const response = await fetch('/api/listings', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(listingData)
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to create listing');
    }
    
    const listing = await response.json();
    
    // Redirect to the new listing
    window.location.href = `/listing/${listing.id}`;
  } catch (error) {
    console.error('Error creating listing:', error);
    
    // Display error message
    const errorElement = form.querySelector('.form-error') || document.createElement('div');
    errorElement.className = 'form-error';
    errorElement.textContent = error.message || 'Failed to create listing. Please try again.';
    
    if (!form.querySelector('.form-error')) {
      form.appendChild(errorElement);
    }
  }
}

/**
 * Handle image upload
 * @param {Event} event - Change event from file input
 */
function handleImageUpload(event) {
  const fileInput = event.target;
  const previewContainer = document.getElementById('image-preview');
  
  if (!previewContainer) return;
  
  // Clear previous preview
  previewContainer.innerHTML = '';
  
  if (fileInput.files && fileInput.files[0]) {
    const reader = new FileReader();
    
    reader.onload = function(e) {
      const img = document.createElement('img');
      img.src = e.target.result;
      img.className = 'preview-image';
      previewContainer.appendChild(img);
      
      // Store the data URL in a hidden input
      const hiddenInput = document.getElementById('image-data') || document.createElement('input');
      hiddenInput.type = 'hidden';
      hiddenInput.name = 'images';
      hiddenInput.id = 'image-data';
      hiddenInput.value = e.target.result;
      
      if (!document.getElementById('image-data')) {
        fileInput.parentNode.appendChild(hiddenInput);
      }
    };
    
    reader.readAsDataURL(fileInput.files[0]);
  }
}

/**
 * Fetch and display a single listing
 * @param {string} listingId - ID of the listing to display
 */
async function fetchListing(listingId) {
  try {
    const response = await fetch(`/api/listings/${listingId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch listing');
    }
    
    const listing = await response.json();
    
    // Display listing details
    displayListingDetails(listing);
    
    return listing;
  } catch (error) {
    console.error('Error fetching listing:', error);
    displayError('Failed to load listing. It may have been removed or is unavailable.');
    return null;
  }
}

/**
 * Display detailed listing information
 * @param {Object} listing - Listing data with user info
 */
function displayListingDetails(listing) {
  const container = document.getElementById('listing-container');
  if (!container) return;
  
  // Clear existing content
  container.innerHTML = '';
  
  // Default image if none provided
  const imageUrl = listing.images && listing.images.length > 0 
    ? listing.images[0] 
    : 'https://images.unsplash.com/photo-1492282442770-077ea21f0f81';
  
  // Create listing detail structure
  const detailElement = createElement('div', { className: 'listing-detail' }, [
    // Images section
    createElement('div', { className: 'listing-images' }, [
      createElement('img', { 
        className: 'listing-image', 
        src: imageUrl,
        alt: listing.title
      })
    ]),
    
    // Info section
    createElement('div', { className: 'listing-info' }, [
      createElement('h1', { className: 'mb-2' }, listing.title),
      createElement('div', { className: 'listing-price' }, formatCurrency(listing.price)),
      
      createElement('div', { className: 'listing-meta' }, [
        createElement('span', { className: 'card-badge' }, listing.type),
        createElement('span', { className: 'card-badge' }, listing.plantType)
      ]),
      
      createElement('p', { className: 'mb-3' }, listing.description),
      
      createElement('div', { className: 'mb-3' }, [
        createElement('strong', {}, 'Location: '),
        createElement('span', {}, listing.location)
      ]),
      
      createElement('div', { className: 'mb-3' }, [
        createElement('strong', {}, 'Will trade for: '),
        createElement('span', {}, listing.tradeFor || 'N/A')
      ]),
      
      createElement('div', { className: 'listing-seller' }, [
        createElement('img', {
          className: 'seller-avatar',
          src: listing.user.profilePic || 'https://images.unsplash.com/photo-1438109382753-8368e7e1e7cf',
          alt: listing.user.name
        }),
        createElement('div', {}, [
          createElement('strong', {}, listing.user.name),
          createElement('p', {}, `Member since ${formatDate(listing.user.createdAt)}`)
        ])
      ]),
      
      createElement('div', { className: 'listing-actions' }, [
        createElement('button', {
          className: 'btn btn-primary btn-lg',
          id: 'contact-seller-btn',
          onclick: () => {
            // Navigate to messages with this user and listing pre-selected
            window.location.href = `/messages?userId=${listing.user.id}&listingId=${listing.id}`;
          }
        }, 'Contact Seller'),
        
        createElement('button', {
          className: 'btn btn-outline btn-lg ml-3',
          id: 'favorite-btn',
          onclick: async (event) => {
            const isFavorite = event.target.dataset.favorite === 'true';
            const newStatus = await toggleFavorite(listing.id, isFavorite);
            event.target.dataset.favorite = newStatus.toString();
            event.target.innerHTML = newStatus ? '★ Favorited' : '☆ Add to Favorites';
            event.target.className = newStatus ? 'btn btn-accent btn-lg ml-3' : 'btn btn-outline btn-lg ml-3';
          },
          dataset: { favorite: 'false' }
        }, '☆ Add to Favorites')
      ])
    ])
  ]);
  
  container.appendChild(detailElement);
  
  // Check if listing is in favorites
  checkFavoriteStatus(listing.id);
}

/**
 * Fetch user's favorite listings
 */
async function fetchFavorites() {
  try {
    const user = await checkAuth();
    if (!user) return;
    
    const response = await fetch('/api/favorites');
    if (!response.ok) {
      throw new Error('Failed to fetch favorites');
    }
    
    const favoriteListings = await response.json();
    favorites = favoriteListings.map(listing => listing.id);
    
    // Update UI for any favorite buttons
    updateFavoriteButtons();
    
    return favoriteListings;
  } catch (error) {
    console.error('Error fetching favorites:', error);
    return [];
  }
}

/**
 * Check if a listing is in the user's favorites and update UI
 * @param {string} listingId - ID of the listing to check
 */
function checkFavoriteStatus(listingId) {
  const favoriteBtn = document.getElementById('favorite-btn');
  if (!favoriteBtn) return;
  
  const isFavorite = favorites.includes(listingId);
  favoriteBtn.dataset.favorite = isFavorite.toString();
  favoriteBtn.innerHTML = isFavorite ? '★ Favorited' : '☆ Add to Favorites';
  favoriteBtn.className = isFavorite ? 'btn btn-accent btn-lg ml-3' : 'btn btn-outline btn-lg ml-3';
}

/**
 * Update all favorite buttons on the page based on favorites list
 */
function updateFavoriteButtons() {
  const favoriteButtons = document.querySelectorAll('[data-listing-id]');
  
  favoriteButtons.forEach(button => {
    const listingId = button.dataset.listingId;
    const isFavorite = favorites.includes(listingId);
    
    button.dataset.favorite = isFavorite.toString();
    button.innerHTML = isFavorite ? '★' : '☆';
    button.className = isFavorite ? 'favorite-btn active' : 'favorite-btn';
  });
}

/**
 * Display an error message
 * @param {string} message - Error message to display
 */
function displayError(message) {
  const container = document.querySelector('.container');
  if (!container) return;
  
  const errorElement = document.createElement('div');
  errorElement.className = 'error-message';
  errorElement.textContent = message;
  
  container.prepend(errorElement);
  
  // Auto-remove after 5 seconds
  setTimeout(() => {
    errorElement.remove();
  }, 5000);
}

// Set up event listeners when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  // Fetch listings for homepage
  const listingsContainer = document.getElementById('listings-container');
  if (listingsContainer) {
    fetchListings();
  }
  
  // Search form
  const searchForm = document.getElementById('search-form');
  if (searchForm) {
    searchForm.addEventListener('submit', handleSearch);
  }
  
  // Filters
  const filterElements = document.querySelectorAll('.filter-select');
  filterElements.forEach(element => {
    element.addEventListener('change', handleFilterChange);
  });
  
  // Create listing form
  const createListingForm = document.getElementById('create-listing-form');
  if (createListingForm) {
    createListingForm.addEventListener('submit', createListing);
  }
  
  // Image upload
  const imageInput = document.getElementById('image-input');
  if (imageInput) {
    imageInput.addEventListener('change', handleImageUpload);
  }
  
  // Single listing page
  const listingContainer = document.getElementById('listing-container');
  const listingId = window.location.pathname.split('/').pop();
  
  if (listingContainer && listingId && listingId !== 'listing') {
    fetchListing(listingId);
  }
  
  // Load user's favorites
  fetchFavorites();
});
