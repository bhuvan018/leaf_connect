// main.js - Main application JavaScript

// DOM elements
const navBtn = document.querySelector('.nav-btn');
const navLinks = document.querySelector('.nav-links');

// Toggle mobile menu
if (navBtn && navLinks) {
  navBtn.addEventListener('click', () => {
    navLinks.classList.toggle('show');
  });
}

// Helper functions
/**
 * Create an element with given properties
 * @param {string} tag - HTML tag name
 * @param {Object} props - Properties to set on the element
 * @param {Array|Node} children - Child nodes to append
 * @returns {HTMLElement} The created element
 */
function createElement(tag, props = {}, children = []) {
  const element = document.createElement(tag);
  
  try {
    // Set properties
    for (const [key, value] of Object.entries(props)) {
      if (key === 'className') {
        element.className = value;
      } else if (key === 'style' && typeof value === 'object') {
        Object.assign(element.style, value);
      } else if (key.startsWith('on') && typeof value === 'function') {
        element.addEventListener(key.substring(2).toLowerCase(), value);
      } else if (key === 'dataset' && typeof value === 'object') {
        // Handle dataset separately
        for (const [dataKey, dataValue] of Object.entries(value)) {
          element.dataset[dataKey] = dataValue;
        }
      } else {
        element[key] = value;
      }
    }
    
    // Append children
    if (Array.isArray(children)) {
      children.forEach(child => {
        if (child !== null && child !== undefined) {
          if (typeof child === 'string') {
            element.appendChild(document.createTextNode(child));
          } else if (child instanceof Node) {
            element.appendChild(child);
          } else {
            console.warn('Invalid child node:', child);
          }
        }
      });
    } else if (children !== null && children !== undefined) {
      if (typeof children === 'string') {
        element.appendChild(document.createTextNode(children));
      } else if (children instanceof Node) {
        element.appendChild(children);
      } else {
        console.warn('Invalid child node:', children);
      }
    }
  } catch (error) {
    console.error('Error creating element:', error, { tag, props, children });
  }
  
  return element;
}

/**
 * Format a date string to a readable format
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted date
 */
function formatDate(dateString) {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }).format(date);
}

/**
 * Format a date string to a time format
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted time
 */
function formatTime(dateString) {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: 'numeric',
    hour12: true
  }).format(date);
}

/**
 * Format currency to display with currency symbol
 * @param {number} amount - Amount to format
 * @returns {string} Formatted currency
 */
function formatCurrency(amount) {
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR'
  }).format(amount);
}

/**
 * Create a card element for a listing
 * @param {Object} listing - Listing data
 * @returns {HTMLElement} Card element
 */
function createListingCard(listing) {
  // Default image if none provided
  const imageUrl = listing.images && listing.images.length > 0 
    ? listing.images[0] 
    : 'https://images.unsplash.com/photo-1492282442770-077ea21f0f81';
  
  return createElement('div', { className: 'card' }, [
    createElement('a', { href: `/listing/${listing.id}` }, [
      createElement('img', { 
        className: 'card-image', 
        src: imageUrl,
        alt: listing.title
      })
    ]),
    createElement('div', { className: 'card-content' }, [
      createElement('h3', { className: 'card-title' }, [
        createElement('a', { href: `/listing/${listing.id}` }, listing.title)
      ]),
      createElement('p', { className: 'card-text' }, listing.description.substring(0, 100) + (listing.description.length > 100 ? '...' : '')),
      createElement('div', { className: 'card-meta' }, [
        createElement('div', {}, [
          createElement('span', { className: 'card-badge' }, listing.type),
          createElement('span', { className: 'card-badge' }, listing.plantType)
        ]),
        createElement('div', { className: 'text-primary' }, formatCurrency(listing.price))
      ]),
      createElement('div', { className: 'card-meta mt-2' }, [
        createElement('div', {}, listing.location),
        createElement('div', {}, formatDate(listing.createdAt))
      ])
    ])
  ]);
}

/**
 * Check if user is authenticated
 * @returns {Promise<Object>} User data if authenticated, null otherwise
 */
async function checkAuth() {
  try {
    const response = await fetch('/api/check-auth');
    const data = await response.json();
    
    // Update UI based on authentication status
    updateAuthUI(data.authenticated, data.user);
    
    return data.authenticated ? data.user : null;
  } catch (error) {
    console.error('Error checking authentication:', error);
    updateAuthUI(false);
    return null;
  }
}

/**
 * Update UI elements based on authentication status
 * @param {boolean} isAuthenticated - Whether user is authenticated
 * @param {Object} user - User data if authenticated
 */
function updateAuthUI(isAuthenticated, user = null) {
  const authLinks = document.querySelectorAll('.auth-link');
  const userLinks = document.querySelectorAll('.user-link');
  const userNameElements = document.querySelectorAll('.user-name');
  
  if (isAuthenticated && user) {
    // Show user-specific links, hide auth links
    authLinks.forEach(el => el.style.display = 'none');
    userLinks.forEach(el => el.style.display = 'block');
    
    // Update user name in UI
    userNameElements.forEach(el => el.textContent = user.name || user.username);
  } else {
    // Show auth links, hide user-specific links
    authLinks.forEach(el => el.style.display = 'block');
    userLinks.forEach(el => el.style.display = 'none');
  }
}

/**
 * Toggle favorite status for a listing
 * @param {string} listingId - ID of the listing
 * @param {boolean} isFavorite - Current favorite status
 * @returns {Promise<boolean>} New favorite status
 */
async function toggleFavorite(listingId, isFavorite) {
  try {
    const response = await fetch('/api/favorites', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        listingId,
        action: isFavorite ? 'remove' : 'add'
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to update favorite');
    }
    
    const data = await response.json();
    return !isFavorite;
  } catch (error) {
    console.error('Error toggling favorite:', error);
    return isFavorite;
  }
}

/**
 * Handle form submission with fetch API
 * @param {Event} event - Form submit event
 * @param {Function} successCallback - Callback on success
 */
function handleFormSubmit(event, successCallback) {
  event.preventDefault();
  
  const form = event.target;
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());
  
  fetch(form.action, {
    method: form.method,
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  .then(response => {
    if (!response.ok) {
      return response.json().then(err => { throw err; });
    }
    return response.json();
  })
  .then(result => {
    if (successCallback) {
      successCallback(result);
    }
  })
  .catch(error => {
    console.error('Form submission error:', error);
    
    // Handle form errors
    const errorElement = form.querySelector('.form-error') || createElement('div', { className: 'form-error' });
    errorElement.textContent = error.message || 'An error occurred. Please try again.';
    
    if (!form.querySelector('.form-error')) {
      form.appendChild(errorElement);
    }
  });
}

// Initialize authentication check on page load
document.addEventListener('DOMContentLoaded', () => {
  checkAuth();
});
