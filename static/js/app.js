import { api } from './api.js';
import { router } from './router.js';
import { renderLoginPage } from './auth.js';
import { renderFeed } from './Pages/feed.js';
import { renderPost } from './Pages/post.js';

function authGuard() {
    return window.__isAuthenticated === true;
}

async function init() {
    try {
        const res = await api.me();
        window.__user = res.data;
        window.__isAuthenticated = true;
    } catch {
        window.__isAuthenticated = false;
    }

    router.addRoute('login', (app) => {
        renderLoginPage(app);
    });

    router.addRoute('register', (app) => {
        renderLoginPage(app);
    });

    router.addRoute('feed', (app) => {
        renderFeed(app);
    }, authGuard);

    router.addRoute('post/:id', (app, params) => {
        renderPost(app, params);
    }, authGuard);

    router.addRoute('create', (app) => {
        app.innerHTML = '<div style="color:white;text-align:center;padding:40px;"><h1>Create Post - Coming Soon</h1></div>';
    }, authGuard);

    router.addRoute('chat', (app) => {
        app.innerHTML = '<div style="color:white;text-align:center;padding:40px;"><h1>Chat - Coming Soon</h1></div>';
    }, authGuard);

    router.addRoute('chat/:userId', (app, params) => {
        app.innerHTML = `<div style="color:white;text-align:center;padding:40px;"><h1>Chat with ${params.userId} - Coming Soon</h1></div>`;
    }, authGuard);

    router.addRoute('notifications', (app) => {
        app.innerHTML = '<div style="color:white;text-align:center;padding:40px;"><h1>Notifications - Coming Soon</h1></div>';
    }, authGuard);

    router.addRoute('profile', (app) => {
        app.innerHTML = '<div style="color:white;text-align:center;padding:40px;"><h1>Profile - Coming Soon</h1></div>';
    }, authGuard);

    router.start();
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}
