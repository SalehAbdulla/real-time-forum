import { api } from './api.js';
import { router } from './router.js';
import { renderLoginPage } from './auth.js';
import { renderFeed } from './pages/feed.js';
import { renderPost } from './pages/post.js';
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
    } catch {
        window.__isAuthenticated = false;
    }

    router.addRoute('login', (app) => {
        renderLoginPage(app);
    }, guestGuard);

    router.addRoute('register', (app) => {
        renderLoginPage(app);
    }, guestGuard);

    router.addRoute('feed', (app) => {
        renderFeed(app);
    }, authGuard);

    router.addRoute('post/:id', (app, params) => {
        renderPost(app, params);
    }, authGuard);

    router.addRoute('create', (app) => {
        app.innerHTML = '<div class="empty-state" style="margin-top:40vh;"><h1 style="font-size:20px;font-weight:500;color:var(--text-secondary);margin-bottom:8px;">Create Post</h1><p style="color:var(--text-muted);">Coming soon</p></div>';
    }, authGuard);

    router.addRoute('chat', (app) => {
        app.innerHTML = '<div class="empty-state" style="margin-top:40vh;"><h1 style="font-size:20px;font-weight:500;color:var(--text-secondary);margin-bottom:8px;">Chat</h1><p style="color:var(--text-muted);">Coming soon</p></div>';
    }, authGuard);

    router.addRoute('chat/:userId', (app, params) => {
        app.innerHTML = `<div class="empty-state" style="margin-top:40vh;"><h1 style="font-size:20px;font-weight:500;color:var(--text-secondary);margin-bottom:8px;">Chat</h1><p style="color:var(--text-muted);">Coming soon</p></div>`;
    }, authGuard);

    router.addRoute('notifications', (app) => {
        app.innerHTML = '<div class="empty-state" style="margin-top:40vh;"><h1 style="font-size:20px;font-weight:500;color:var(--text-secondary);margin-bottom:8px;">Notifications</h1><p style="color:var(--text-muted);">Coming soon</p></div>';
    }, authGuard);

    router.addRoute('profile', (app) => {
        app.innerHTML = '<div class="empty-state" style="margin-top:40vh;"><h1 style="font-size:20px;font-weight:500;color:var(--text-secondary);margin-bottom:8px;">Profile</h1><p style="color:var(--text-muted);">Coming soon</p></div>';
    }, authGuard);

    router.start();
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}