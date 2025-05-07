// profile.js - Profile-related functionality

/**
 * Fetch user profile data
 * @param {string} userId - User ID to fetch, or null for current user
 */
async function fetchUserProfile(userId = null) {
  try {
    // Check authentication first
    const currentUser = await checkAuth();
    if (!currentUser && !userId) {
      window.location.href = '/login';
      return null;
    }

    // Determine which endpoint to use
    const url = userId ? `/api/users/${userId}` : '/api/users/current';
    
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('Failed to fetch profile');
    }
    
    const profile = await response.json();
    
    // Display the profile information
    displayProfile(profile);
    
    // Load user's listings
    fetchUserListings(profile.id);
    
    return profile;
  } catch (error) {
    console.error('Error fetching profile:', error);
    displayError('Failed to load profile. Please try again later.');
    return null;
  }
}

/**
 * Display user profile information
 * @param {Object} profile - User profile data
 */
function displayProfile(profile) {
  // Profile name
  const nameElements = document.querySelectorAll('.profile-name');
  nameElements.forEach(el => {
    el.textContent = profile.name || profile.username;
  });
  
  // Username
  const usernameElements = document.querySelectorAll('.profile-username');
  usernameElements.forEach(el => {
    el.textContent = profile.username;
  });
  
  // Location
  const locationElements = document.querySelectorAll('.profile-location');
  locationElements.forEach(el => {
    el.textContent = profile.location || 'Location not specified';
  });
  
  // Bio
  const bioElements = document.querySelectorAll('.profile-bio');
  bioElements.forEach(el => {
    el.textContent = profile.bio || 'No bio provided';
  });
  
  // Profile image
  const avatarElements = document.querySelectorAll('.profile-avatar');
  avatarElements.forEach(el => {
    el.src = profile.profilePic || 'https://images.unsplash.com/photo-1438109382753-8368e7e1e7cf';
    el.alt = profile.name || profile.username;
  });
  
  // Join date
  const joinDateElements = document.querySelectorAll('.profile-join-date');
  joinDateElements.forEach(el => {
    el.textContent = formatDate(profile.createdAt);
  });
  
  // Check if this is the current user's profile
  checkAuth().then(currentUser => {
    if (currentUser && currentUser.id === profile.id) {
      // Show edit profile button
      const editButtons = document.querySelectorAll('.edit-profile-btn');
      editButtons.forEach(btn => {
        btn.style.display = 'block';
      });
    }
  });
}

/**
 * Fetch and display user's listings
 * @param {string} userId - User ID to fetch listings for
 */
async function fetchUserListings(userId) {
  try {
    const response = await fetch(`/api/listings?userId=${userId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch user listings');
    }
    
    const listings = await response.json();
    
    // Display the listings
    const container = document.getElementById('user-listings');
    if (!container) return;
    
    // Clear existing content
    container.innerHTML = '';
    
    if (listings.length === 0) {
      container.innerHTML = '<p class="text-center">This user hasn\'t created any listings yet.</p>';
      return;
    }
    
    // Create grid for listings
    const grid = createElement('div', { className: 'grid' });
    
    // Add listing cards
    listings.forEach(listing => {
      grid.appendChild(createListingCard(listing));
    });
    
    container.appendChild(grid);
    
  } catch (error) {
    console.error('Error fetching user listings:', error);
    const container = document.getElementById('user-listings');
    if (container) {
      container.innerHTML = '<p class="text-center text-error">Failed to load listings. Please try again later.</p>';
    }
  }
}

/**
 * Handle profile update form submission
 * @param {Event} event - Form submit event
 */
async function updateProfile(event) {
  event.preventDefault();
  
  // Get current user
  const currentUser = await checkAuth();
  if (!currentUser) {
    window.location.href = '/login';
    return;
  }
  
  const form = event.target;
  const formData = new FormData(form);
  const profileData = {};
  
  // Build profile data object
  formData.forEach((value, key) => {
    profileData[key] = value;
  });
  
  try {
    const response = await fetch(`/api/users/${currentUser.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(profileData)
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to update profile');
    }
    
    const updatedProfile = await response.json();
    
    // Show success message
    const successMessage = document.createElement('div');
    successMessage.className = 'alert alert-success';
    successMessage.textContent = 'Profile updated successfully!';
    form.prepend(successMessage);
    
    // Remove message after 3 seconds
    setTimeout(() => {
      successMessage.remove();
    }, 3000);
    
    // Update displayed profile information
    displayProfile(updatedProfile);
    
  } catch (error) {
    console.error('Error updating profile:', error);
    
    // Display error message
    const errorElement = form.querySelector('.form-error') || document.createElement('div');
    errorElement.className = 'form-error';
    errorElement.textContent = error.message || 'Failed to update profile. Please try again.';
    
    if (!form.querySelector('.form-error')) {
      form.prepend(errorElement);
    }
  }
}

/**
 * Handle profile image upload
 * @param {Event} event - Change event from file input
 */
function handleProfileImageUpload(event) {
  const fileInput = event.target;
  const previewElement = document.getElementById('profile-image-preview');
  
  if (!previewElement) return;
  
  if (fileInput.files && fileInput.files[0]) {
    const reader = new FileReader();
    
    reader.onload = function(e) {
      previewElement.src = e.target.result;
      
      // Store the data URL in a hidden input
      const hiddenInput = document.getElementById('profile-image-data') || document.createElement('input');
      hiddenInput.type = 'hidden';
      hiddenInput.name = 'profilePic';
      hiddenInput.id = 'profile-image-data';
      hiddenInput.value = e.target.result;
      
      if (!document.getElementById('profile-image-data')) {
        fileInput.parentNode.appendChild(hiddenInput);
      }
    };
    
    reader.readAsDataURL(fileInput.files[0]);
  }
}

/**
 * Toggle between view and edit mode for profile
 */
function toggleEditProfile() {
  const viewSection = document.getElementById('profile-view');
  const editSection = document.getElementById('profile-edit');
  
  if (!viewSection || !editSection) return;
  
  // Toggle visibility
  if (viewSection.style.display === 'none') {
    viewSection.style.display = 'block';
    editSection.style.display = 'none';
  } else {
    viewSection.style.display = 'none';
    editSection.style.display = 'block';
    
    // Populate edit form with current values
    checkAuth().then(currentUser => {
      if (!currentUser) return;
      
      const nameInput = document.getElementById('edit-name');
      const locationInput = document.getElementById('edit-location');
      const bioInput = document.getElementById('edit-bio');
      
      if (nameInput) nameInput.value = currentUser.name || '';
      if (locationInput) locationInput.value = currentUser.location || '';
      if (bioInput) bioInput.value = currentUser.bio || '';
    });
  }
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
  // Check if on profile page
  const profileContainer = document.getElementById('profile-container');
  if (profileContainer) {
    // Extract user ID from URL if present
    const pathParts = window.location.pathname.split('/');
    const userId = pathParts[pathParts.length - 1] !== 'profile' ? pathParts[pathParts.length - 1] : null;
    
    // Load profile
    fetchUserProfile(userId);
  }
  
  // Edit profile form
  const editProfileForm = document.getElementById('edit-profile-form');
  if (editProfileForm) {
    editProfileForm.addEventListener('submit', updateProfile);
  }
  
  // Profile image upload
  const profileImageInput = document.getElementById('profile-image-input');
  if (profileImageInput) {
    profileImageInput.addEventListener('change', handleProfileImageUpload);
  }
  
  // Edit profile button
  const editProfileBtn = document.getElementById('edit-profile-btn');
  if (editProfileBtn) {
    editProfileBtn.addEventListener('click', toggleEditProfile);
  }
  
  // Cancel edit button
  const cancelEditBtn = document.getElementById('cancel-edit-btn');
  if (cancelEditBtn) {
    cancelEditBtn.addEventListener('click', toggleEditProfile);
  }
});
