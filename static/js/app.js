import { api } from './api.js';
import { router } from './router.js';
import { renderLoginPage } from './auth.js';
import { renderFeed } from './pages/feed.js';
import { renderPost } from './pages/post.js';
import { renderCreatePost } from './pages/createPost.js';
import { renderProfile } from './pages/profile.js';
import { renderChat } from './pages/chat.js';
import { renderNotifications } from './pages/notifications.js';
import { ws } from './websocket.js';
import { createParticles } from './utils.js';

function authGuard() {
    return window.__isAuthenticated === true;
}

function guestGuard() {
    if (window.__isAuthenticated === true) {
        router.navigate('feed');
        return false;
    }
    return true;
}

async function init() {
    // Initialize ambient particles
    createParticles();

    try {
        const res = await api.me();
        window.__user = res.data;
        window.__isAuthenticated = true;
        // Connect WebSocket — the session cookie from the response is HttpOnly,
        // so the browser will automatically send it on the WS upgrade request.
        ws.connect();
    } catch {
        window.__isAuthenticated = false;
    }

    router.addRoute('login', (app) => {
        renderLoginPage(app);
    }, guestGuard);

    router.addRoute('register', (app) => {
        renderLoginPage(app);
    }, guestGuard);

    router.addRoute('feed', (app, params, queryString) => {
        renderFeed(app, params, queryString);
    }, authGuard);

    router.addRoute('post/:id', (app, params) => {
        renderPost(app, params);
    }, authGuard);

    router.addRoute('create', (app) => {
        renderCreatePost(app);
    }, authGuard);

    router.addRoute('chat', (app, params, queryString) => {
        renderChat(app, params, queryString);
    }, authGuard);

    router.addRoute('chat/:userId', (app, params, queryString) => {
        renderChat(app, params, queryString);
    }, authGuard);

    router.addRoute('notifications', (app) => {
        renderNotifications(app);
    }, authGuard);

    router.addRoute('profile', (app) => {
        renderProfile(app);
    }, authGuard);

    router.start();

    // Global: fetch unread notification count and listen for real-time notifications
    if (window.__isAuthenticated) {
        initNotificationBadge();
    }
}

/**
 * Fetches the initial unread notification count and sets up the WebSocket listener
 * to update the sidebar badge across all pages in real time.
 */
async function initNotificationBadge() {
    try {
        const res = await api.getUnreadCount();
        const count = res.data?.count || 0;
        updateGlobalBadge(count);
    } catch (err) {
        console.error('Failed to load unread count:', err);
    }

    // Listen for real-time notification events via WebSocket
    if (!window._globalNotifCleanup) {
        window._globalNotifCleanup = ws.on('notification', (payload) => {
            if (!payload) return;
            // Increment the global unread count
            window.__unreadNotifCount = (window.__unreadNotifCount || 0) + 1;
            updateGlobalBadge(window.__unreadNotifCount);
        });
    }
}

/**
 * Updates the sidebar notification badge element.
 */
function updateGlobalBadge(count) {
    window.__unreadNotifCount = count;
    const notifNavItem = document.querySelector('.sidebar-nav-item[data-route="notifications"]');
    if (!notifNavItem) return;

    let badgeEl = notifNavItem.querySelector('.nav-badge');
    if (count > 0) {
        if (!badgeEl) {
            badgeEl = document.createElement('span');
            badgeEl.className = 'nav-badge';
            notifNavItem.appendChild(badgeEl);
        }
        badgeEl.textContent = count > 99 ? '99+' : count;
        badgeEl.style.display = '';
    } else if (badgeEl) {
        badgeEl.style.display = 'none';
    }
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}