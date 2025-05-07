// auth.js - Authentication-related functionality

/**
 * Handle user registration
 * @param {Event} event - Form submit event
 */
function handleRegister(event) {
  event.preventDefault();
  
  const form = event.target;
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());
  
  // Debug info
  console.log('Form data:', data);
  
  // Validate password match
  if (data.password !== data.confirmPassword) {
    showError(form, 'Passwords do not match');
    return;
  }
  
  // Remove confirmPassword from data
  delete data.confirmPassword;
  
  // Debug the JSON being sent
  console.log('JSON to send:', JSON.stringify(data));
  
  fetch('/api/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  .then(response => {
    if (!response.ok) {
      return response.text().then(text => {
        try {
          const err = JSON.parse(text);
          throw new Error(err.message || 'Registration failed');
        } catch (e) {
          console.error('Error parsing response:', text);
          throw new Error(text || 'Registration failed');
        }
      });
    }
    return response.json();
  })
  .then(user => {
    // Registration successful - redirect to dashboard
    window.location.href = '/dashboard';
  })
  .catch(error => {
    showError(form, error.message);
  });
}

/**
 * Handle user login
 * @param {Event} event - Form submit event
 */
function handleLogin(event) {
  event.preventDefault();
  
  const form = event.target;
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());
  
  fetch('/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  .then(response => {
    if (!response.ok) {
      return response.text().then(text => {
        try {
          const err = JSON.parse(text);
          throw new Error(err.message || 'Invalid email or password');
        } catch (e) {
          console.error('Error parsing response:', text);
          throw new Error(text || 'Login failed');
        }
      });
    }
    return response.json();
  })
  .then(user => {
    // Login successful - redirect to dashboard
    window.location.href = '/dashboard';
  })
  .catch(error => {
    showError(form, error.message);
  });
}

/**
 * Handle user logout
 */
function handleLogout() {
  fetch('/api/logout', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    }
  })
  .then(response => response.json())
  .then(data => {
    // Logout successful - redirect to home page
    window.location.href = '/';
  })
  .catch(error => {
    console.error('Logout error:', error);
  });
}

/**
 * Display error message in form
 * @param {HTMLFormElement} form - Form element
 * @param {string} message - Error message
 */
function showError(form, message) {
  // Find or create error element
  let errorElement = form.querySelector('.form-error');
  
  if (!errorElement) {
    errorElement = document.createElement('div');
    errorElement.className = 'form-error';
    form.appendChild(errorElement);
  }
  
  errorElement.textContent = message;
  errorElement.style.display = 'block';
}

// Set up event listeners when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  // Register form
  const registerForm = document.getElementById('register-form');
  if (registerForm) {
    registerForm.addEventListener('submit', handleRegister);
  }
  
  // Login form
  const loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', handleLogin);
  }
  
  // Logout button
  const logoutBtn = document.getElementById('logout-btn');
  if (logoutBtn) {
    logoutBtn.addEventListener('click', (event) => {
      event.preventDefault();
      handleLogout();
    });
  }
});
