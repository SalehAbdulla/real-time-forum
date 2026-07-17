/**
 * FeedLayout — wraps page content with header, users, categories
 */
import { createFeedHeader } from '../components/FeedHeader.js';
import { createUsersScroll } from '../components/UsersScroll.js';
import { createCategoryPills } from '../components/CategoryPills.js';

export function createFeedLayout(contentArea, options = {}) {
    const {
        showUsers = true,
        showCategories = true,
        onCategoryChange = null,
    } = options;
    
    const container = document.createElement('div');
    container.className = 'feed-container';
    
    // Header
    container.appendChild(createFeedHeader());
    
    // Users scroll
    if (showUsers) {
        container.appendChild(createUsersScroll());
    }
    
    // Categories
    if (showCategories) {
        container.appendChild(createCategoryPills('all', onCategoryChange));
    }
    
    // Content area (posts feed, etc.)
    if (typeof contentArea === 'function') {
        const area = document.createElement('div');
        area.className = 'posts-feed';
        area.id = 'posts-feed';
        contentArea(area);
        container.appendChild(area);
    } else if (contentArea instanceof HTMLElement) {
        container.appendChild(contentArea);
    }
    
    return container;
}
