// messages.js - Message-related functionality

// State
let currentConversations = [];
let currentConversation = null;
let currentUser = null;

/**
 * Fetch all conversations
 */
async function fetchConversations() {
  try {
    console.log('Starting fetchConversations');
    
    // Get current user
    const user = await checkAuth();
    console.log('User authentication check:', user ? 'Authenticated' : 'Not authenticated');
    
    if (!user) {
      console.log('No authenticated user, redirecting to login');
      window.location.href = '/login';
      return;
    }
    
    currentUser = user;
    console.log('Current user set:', currentUser);
    
    console.log('Fetching conversations from API');
    const response = await fetch('/api/conversations');
    console.log('API response status:', response.status);
    
    if (!response.ok) {
      throw new Error('Failed to fetch conversations');
    }
    
    const conversations = await response.json();
    console.log('Received conversations:', conversations);
    currentConversations = conversations;
    
    // Display the conversations
    console.log('Displaying conversations');
    displayConversations(conversations);
    
    // Check URL for specific conversation
    const urlParams = new URLSearchParams(window.location.search);
    const userId = urlParams.get('userId');
    console.log('URL userId parameter:', userId);
    
    if (userId) {
      // Open conversation with specified user
      console.log('Opening conversation with specified user:', userId);
      fetchConversation(userId);
    } else if (conversations.length > 0) {
      // Open first conversation by default
      console.log('Opening first conversation by default:', conversations[0].userId);
      fetchConversation(conversations[0].userId);
    } else {
      // Show empty state if no conversations
      console.log('No conversations, showing empty state');
      displayEmptyConversation();
    }
    
    return conversations;
  } catch (error) {
    console.error('Error fetching conversations:', error);
    displayError('Failed to load conversations. Please try again later.');
    return [];
  }
}

/**
 * Display conversations list
 * @param {Array} conversations - Array of conversation data
 */
function displayConversations(conversations) {
  console.log('Starting displayConversations with', conversations.length, 'items');
  
  const container = document.getElementById('conversations-list');
  if (!container) {
    console.error('Conversations list container not found');
    return;
  }
  
  // Clear existing content
  container.innerHTML = '';
  console.log('Cleared existing content in conversations list');
  
  if (conversations.length === 0) {
    console.log('No conversations to display, showing empty state');
    container.innerHTML = '<p class="text-center p-3">No conversations yet.</p>';
    return;
  }
  
  // Add conversation items
  console.log('Adding conversation items to list');
  conversations.forEach((conversation, index) => {
    try {
      console.log(`Creating conversation item ${index + 1}/${conversations.length}:`, conversation.userId);
      
      // Fix potential issues with lastActivity not being a valid date
      let formattedDate = '';
      try {
        if (conversation.lastActivity) {
          formattedDate = formatDate(conversation.lastActivity);
        }
      } catch (dateError) {
        console.warn('Error formatting date:', dateError);
        formattedDate = 'Recent';
      }
      
      // Create elements with proper error handling
      const header = createElement('div', { className: 'conversation-header' }, [
        createElement('img', {
          className: 'conversation-avatar',
          src: conversation.profilePic || 'https://images.unsplash.com/photo-1438109382753-8368e7e1e7cf',
          alt: conversation.username || 'User'
        }),
        createElement('div', { className: 'conversation-info' }, [
          createElement('div', { className: 'conversation-username' }, conversation.username || 'Unknown User'),
          createElement('div', { className: 'conversation-time' }, formattedDate)
        ])
      ]);
      
      const preview = createElement('div', { className: 'conversation-preview' }, 
        conversation.lastMessage || 'Start a conversation'
      );
      
      // Only create badge if unread > 0
      let badge = null;
      if (conversation.unread && conversation.unread > 0) {
        badge = createElement('div', { className: 'conversation-badge' }, String(conversation.unread));
      }
      
      // Create children array without null elements
      const childrenArray = [header, preview];
      if (badge) childrenArray.push(badge);
      
      // Create the item with the valid children
      const item = createElement('div', {
        className: 'conversation-item',
        dataset: { userId: conversation.userId },
        onclick: () => fetchConversation(conversation.userId)
      }, childrenArray);
      
      console.log(`Appending conversation item ${index + 1} to container`);
      container.appendChild(item);
    } catch (err) {
      console.error(`Error creating conversation item ${index}:`, err, conversation);
    }
  });
  
  console.log('Finished displaying conversations');
}

/**
 * Fetch messages for a specific conversation
 * @param {string} userId - ID of the user to fetch conversation with
 */
async function fetchConversation(userId) {
  try {
    console.log('Starting fetchConversation for userId:', userId);
    
    console.log('Fetching conversation from API');
    const response = await fetch(`/api/conversations/${userId}`);
    console.log('API response status:', response.status);
    
    if (!response.ok) {
      throw new Error('Failed to fetch conversation');
    }
    
    const conversation = await response.json();
    console.log('Received conversation data:', conversation);
    currentConversation = conversation;
    
    // Display the messages
    console.log('Displaying conversation');
    displayConversation(conversation);
    
    // Update active state in conversation list
    console.log('Updating active conversation in list');
    updateActiveConversation(userId);
    
    return conversation;
  } catch (error) {
    console.error('Error fetching conversation:', error);
    displayError('Failed to load messages. Please try again later.');
    return null;
  }
}

/**
 * Display messages for a specific conversation
 * @param {Object} conversation - Conversation data
 */
function displayConversation(conversation) {
  console.log('Starting displayConversation', conversation);
  
  const headerContainer = document.getElementById('message-header');
  const bodyContainer = document.getElementById('message-body');
  const formContainer = document.getElementById('message-form');
  
  console.log('Containers found:', {
    header: !!headerContainer,
    body: !!bodyContainer,
    form: !!formContainer
  });
  
  if (!headerContainer || !bodyContainer || !formContainer) {
    console.error('Missing required DOM containers for conversation display');
    return;
  }
  
  try {
    // Set up header
    console.log('Setting up conversation header');
    headerContainer.innerHTML = '';
    const headerContent = createElement('div', { className: 'message-header-content' }, [
      createElement('img', {
        className: 'conversation-avatar',
        src: conversation.profilePic || 'https://images.unsplash.com/photo-1438109382753-8368e7e1e7cf',
        alt: conversation.username
      }),
      createElement('div', { className: 'message-user-info' }, [
        createElement('div', { className: 'conversation-username' }, conversation.username)
      ])
    ]);
    headerContainer.appendChild(headerContent);
    
    // Set up message body
    console.log('Setting up message body');
    bodyContainer.innerHTML = '';
    
    if (!conversation.messages || conversation.messages.length === 0) {
      console.log('No messages in conversation, showing empty state');
      bodyContainer.innerHTML = '<p class="text-center p-3">No messages yet. Send a message to start the conversation!</p>';
    } else {
      // Display messages
      console.log(`Displaying ${conversation.messages.length} messages`);
      conversation.messages.forEach((message, index) => {
        console.log(`Processing message ${index + 1}/${conversation.messages.length}`);
        try {
          const isSentByCurrentUser = message.fromUser.id === currentUser.id;
          
          const messageElement = createElement('div', {
            className: `message-bubble ${isSentByCurrentUser ? 'message-sent' : 'message-received'}`
          }, [
            createElement('div', { className: 'message-content' }, message.content),
            createElement('div', { className: 'message-time' }, formatTime(message.createdAt))
          ]);
          
          bodyContainer.appendChild(messageElement);
        } catch (err) {
          console.error(`Error creating message element ${index}:`, err, message);
        }
      });
      
      // Scroll to bottom
      bodyContainer.scrollTop = bodyContainer.scrollHeight;
    }
    
    // Set up message form
    console.log('Setting up message form');
    formContainer.innerHTML = '';
    const formElement = createElement('form', {
      id: 'send-message-form',
      onsubmit: (event) => sendMessage(event, conversation.userId)
    }, [
      createElement('input', {
        type: 'text',
        name: 'message',
        className: 'form-control message-input',
        placeholder: 'Type your message...',
        required: true
      }),
      createElement('button', {
        type: 'submit',
        className: 'btn btn-primary'
      }, 'Send')
    ]);
    formContainer.appendChild(formElement);
    
    // Check if this is a new conversation from listing
    const urlParams = new URLSearchParams(window.location.search);
    const listingId = urlParams.get('listingId');
    
    if (listingId && conversation.messages.length === 0) {
      console.log('New conversation with listing ID, fetching listing details');
      // Fetch listing details to show in message area
      fetchListingForMessage(listingId);
    }
  } catch (error) {
    console.error('Error in displayConversation:', error);
    displayError('Failed to display conversation. Please try refreshing the page.');
  }
}

/**
 * Display empty conversation state
 */
function displayEmptyConversation() {
  const headerContainer = document.getElementById('message-header');
  const bodyContainer = document.getElementById('message-body');
  const formContainer = document.getElementById('message-form');
  
  if (!headerContainer || !bodyContainer || !formContainer) return;
  
  // Empty header
  headerContainer.innerHTML = '<div class="message-header-content">New Message</div>';
  
  // Empty body with instructions
  bodyContainer.innerHTML = `
    <div class="empty-messages">
      <p>Select a conversation from the list or start a new one from a listing page.</p>
    </div>
  `;
  
  // Empty form
  formContainer.innerHTML = '';
}

/**
 * Update active conversation in the list
 * @param {string} userId - ID of the active conversation user
 */
function updateActiveConversation(userId) {
  const items = document.querySelectorAll('.conversation-item');
  
  items.forEach(item => {
    if (item.dataset.userId === userId) {
      item.classList.add('active');
    } else {
      item.classList.remove('active');
    }
  });
}

/**
 * Send a message
 * @param {Event} event - Form submit event
 * @param {string} recipientId - ID of the recipient
 */
async function sendMessage(event, recipientId) {
  event.preventDefault();
  
  const form = event.target;
  const messageInput = form.querySelector('input[name="message"]');
  
  if (!messageInput || !messageInput.value.trim()) {
    return;
  }
  
  const message = messageInput.value.trim();
  
  // Get listing ID from URL or use the most recent one
  const urlParams = new URLSearchParams(window.location.search);
  let listingId = urlParams.get('listingId');
  
  if (!listingId && currentConversation && currentConversation.messages && currentConversation.messages.length > 0) {
    // Use listing ID from most recent message
    listingId = currentConversation.messages[0].listing.id;
  }
  
  if (!listingId) {
    displayError('Cannot send message: No listing selected.');
    return;
  }
  
  try {
    const response = await fetch('/api/messages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        toID: recipientId,
        listingID: listingId,
        content: message
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to send message');
    }
    
    // Clear input
    messageInput.value = '';
    
    // Reload conversation to show new message
    await fetchConversation(recipientId);
    
    // Also reload conversations list to update last message
    await fetchConversations();
  } catch (error) {
    console.error('Error sending message:', error);
    displayError('Failed to send message. Please try again.');
  }
}

/**
 * Fetch listing details for new conversation
 * @param {string} listingId - ID of the listing
 */
async function fetchListingForMessage(listingId) {
  try {
    const response = await fetch(`/api/listings/${listingId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch listing');
    }
    
    const listing = await response.json();
    
    // Display listing info in message area
    const bodyContainer = document.getElementById('message-body');
    if (!bodyContainer) return;
    
    bodyContainer.innerHTML = '';
    bodyContainer.appendChild(
      createElement('div', { className: 'listing-message-info' }, [
        createElement('p', {}, 'Start the conversation about:'),
        createElement('div', { className: 'listing-message-card' }, [
          createElement('img', {
            className: 'listing-message-image',
            src: listing.images && listing.images.length > 0 ? listing.images[0] : 'https://images.unsplash.com/photo-1492282442770-077ea21f0f81',
            alt: listing.title
          }),
          createElement('div', { className: 'listing-message-details' }, [
            createElement('h4', {}, listing.title),
            createElement('p', {}, formatCurrency(listing.price))
          ])
        ])
      ])
    );
    
  } catch (error) {
    console.error('Error fetching listing for message:', error);
  }
}

/**
 * Display an error message
 * @param {string} message - Error message to display
 */
function displayError(message) {
  console.log('Displaying error:', message);
  const container = document.querySelector('.messages-container');
  if (!container) {
    console.error('Could not find messages-container to display error');
    return;
  }
  
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
  // Load conversations
  const messagesContainer = document.querySelector('.messages-container');
  if (messagesContainer) {
    fetchConversations();
  }
});
