/**
 * UsersScroll component — horizontal scrollable online users list
 */
import { api } from '../api.js';
import { router } from '../router.js';
import { createAvatarWithStatus } from './Avatar.js';
import { escapeHtml } from '../utils.js';

export function createUsersScroll() {
    const wrapper = document.createElement('div');
    wrapper.className = 'users-wrapper';
    
    const scroll = document.createElement('div');
    scroll.className = 'users-scroll';
    scroll.id = 'users-scroll';
    scroll.innerHTML = '<div class="loading-spinner">Loading users...</div>';
    
    wrapper.appendChild(scroll);
    
    // Load users async
    loadUsers(scroll);
    
    return wrapper;
}

async function loadUsers(container) {
    try {
        const res = await api.getChatUsers();
        const users = res.data || [];
        
        if (users.length === 0) {
            container.innerHTML = '';
            return;
        }
        
        container.innerHTML = '';
        
        users.forEach(user => {
            const card = document.createElement('div');
            card.className = 'user-card';
            card.dataset.userId = user.userId;
            
            const avatarWrapper = createAvatarWithStatus(user.nickname, user.isOnline === 1);
            card.appendChild(avatarWrapper);
            
            const label = document.createElement('div');
            label.className = 'user-nickname';
            label.textContent = escapeHtml(user.nickname);
            card.appendChild(label);
            
            card.addEventListener('click', () => {
                router.navigate(`chat/${user.userId}`);
            });
            
            container.appendChild(card);
        });
    } catch (err) {
        container.innerHTML = '';
    }
}