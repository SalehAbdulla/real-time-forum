class Router {
    constructor() {
        this.routes = {};
        this.guards = {};
        this.currentRoute = null;

        window.addEventListener('hashchange', () => this.resolve());
    }

    addRoute(pattern, handler, guard = null) {
        this.routes[pattern] = handler;
        if (guard) {
            this.guards[pattern] = guard;
        }
    }

    navigate(path) {
        window.location.hash = path;
    }

    resolve() {
        const hash = window.location.hash.slice(1) || '/login';
        const [path, ...rest] = hash.split('?');
        const queryString = rest.join('?');
        const params = {};

        for (const pattern in this.routes) {
            const match = this.matchRoute(pattern, path);
            if (match) {
                const guard = this.guards[pattern];
                if (guard && !guard()) {
                    this.navigate('');
                    return;
                }

                this.currentRoute = pattern;
                const app = document.getElementById('app');

                // Auth routes — render standalone without the desktop card frame
                if (path === 'login' || path === 'register') {
                    app.innerHTML = `<div id="app-content"></div>`;
                    const content = document.getElementById('app-content');
                    this.routes[pattern](content, match.params, queryString);
                    return;
                }

                // Authenticated routes — render inside the desktop card frame
                const name = window.__user ? (window.__user.nickname || window.__user.firstName || 'User') : 'User';
                const initials = name.substring(0, 2).toUpperCase();

                app.innerHTML = `
                    <div class="app-frame">
                        <div class="app-inner">
                            <!-- Left Sidebar -->
                            <div class="app-sidebar">
                                <div class="sidebar-logo">
                                    <span>forum</span><span class="logo-dot"></span>
                                </div>
                                <nav class="sidebar-nav">
                                    <button class="sidebar-nav-item ${path === 'feed' ? 'active' : ''}" data-route="feed">
                                        <svg viewBox="0 0 20 20" fill="none">
                                            <path d="M2 10L10 2L18 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                            <path d="M4 8V16H16V8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        </svg>
                                        Home
                                    </button>
                                    <button class="sidebar-nav-item ${path === 'chat' ? 'active' : ''}" data-route="chat">
                                        <svg viewBox="0 0 20 20" fill="none">
                                            <path d="M2 4C2 3.44772 2.44772 3 3 3H17C17.5523 3 18 3.44772 18 4V14C18 14.5523 17.5523 15 17 15H6L3 18V4Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        </svg>
                                        Chat
                                    </button>
                                    <button class="sidebar-nav-item ${path === 'notifications' ? 'active' : ''}" data-route="notifications">
                                        <svg viewBox="0 0 20 20" fill="none">
                                            <path d="M10 2C7.79086 2 6 3.79086 6 6V9C6 9.55228 5.55228 10 5 10H4V12H16V10H15C14.4477 10 14 9.55228 14 9V6C14 3.79086 12.2091 2 10 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                            <path d="M8 14V15C8 16.1046 8.89543 17 10 17C11.1046 17 12 16.1046 12 15V14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        </svg>
                                        Notifications
                                    </button>
                                    <button class="sidebar-nav-item ${path === 'create' ? 'active' : ''}" data-route="create">
                                        <svg viewBox="0 0 20 20" fill="none">
                                            <path d="M10 4V16M4 10H16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        </svg>
                                        Create Post
                                    </button>
                                </nav>

                                <!-- Categories in sidebar (desktop) -->
                                <div class="sidebar-categories">
                                    <span class="sidebar-categories-label">Categories</span>
                                    <div class="sidebar-categories-list">
                                        <button class="sidebar-cat-pill active" data-category="all">All</button>
                                        <button class="sidebar-cat-pill" data-category="1">Tech</button>
                                        <button class="sidebar-cat-pill" data-category="2">General</button>
                                        <button class="sidebar-cat-pill" data-category="3">Dev</button>
                                        <button class="sidebar-cat-pill" data-category="4">Gaming</button>
                                        <button class="sidebar-cat-pill" data-category="5">Q&A</button>
                                        <button class="sidebar-cat-pill" data-category="6">Random</button>
                                        <button class="sidebar-cat-pill" data-category="7">Life</button>
                                        <button class="sidebar-cat-pill" data-category="8">Sport</button>
                                    </div>
                                </div>

                                <div class="sidebar-profile" id="sidebar-profile">
                                    <div class="sidebar-profile-avatar">${initials}</div>
                                    <div class="sidebar-profile-info">
                                        <span class="sidebar-profile-name">${name}</span>
                                        <span class="sidebar-profile-status">Online</span>
                                    </div>
                                </div>
                            </div>

                            <!-- Main Panel -->
                            <div class="app-panel" id="app-panel"></div>

                            <!-- Right Sidebar (Users) -->
                            <div class="app-rightbar">
                                <div class="rightbar-header">
                                    <span>Online</span>
                                </div>
                                <div class="rightbar-users" id="rightbar-users">
                                    <div class="loading-spinner" style="padding:20px 0;">Loading...</div>
                                </div>
                            </div>
                        </div>
                    </div>
                `;

                // Sidebar navigation clicks
                document.querySelectorAll('.sidebar-nav-item').forEach(item => {
                    item.addEventListener('click', () => {
                        const route = item.dataset.route;
                        if (route) this.navigate(route);
                    });
                });

                // Profile click
                const profile = document.getElementById('sidebar-profile');
                if (profile) {
                    profile.addEventListener('click', () => this.navigate('profile'));
                }

                // Sidebar category clicks
                document.querySelectorAll('.sidebar-cat-pill').forEach(pill => {
                    pill.addEventListener('click', () => {
                        document.querySelectorAll('.sidebar-cat-pill').forEach(p => p.classList.remove('active'));
                        pill.classList.add('active');
                        window.dispatchEvent(new CustomEvent('category-change', {
                            detail: { category: pill.dataset.category }
                        }));
                    });
                });

                // Load users into right sidebar
                this.loadRightbarUsers();

                const panel = document.getElementById('app-panel');
                this.routes[pattern](panel, match.params, queryString);
                return;
            }
        }

        this.navigate('login');
    }

    async loadRightbarUsers() {
        const container = document.getElementById('rightbar-users');
        if (!container) return;
        try {
            const res = await fetch('/api/v1/messages/users', { credentials: 'include' });
            const data = await res.json();
            const users = data.data || [];
            if (users.length === 0) {
                container.innerHTML = '<div class="rightbar-empty">No users online</div>';
                return;
            }
            container.innerHTML = users.map(user => `
                <div class="rightbar-user" data-user-id="${user.userId}">
                    <div class="rightbar-user-avatar">${user.nickname ? user.nickname.substring(0, 2).toUpperCase() : '?'}</div>
                    <div class="rightbar-user-info">
                        <span class="rightbar-user-name">${user.nickname || 'User'}</span>
                        <span class="rightbar-user-status ${user.isOnline === 1 ? 'online' : 'offline'}">
                            ${user.isOnline === 1 ? 'Online' : 'Offline'}
                        </span>
                    </div>
                    <div class="rightbar-status-dot ${user.isOnline === 1 ? 'online' : 'offline'}"></div>
                </div>
            `).join('');

            container.querySelectorAll('.rightbar-user').forEach(el => {
                el.addEventListener('click', () => {
                    this.navigate(`chat/${el.dataset.userId}`);
                });
            });
        } catch (err) {
            container.innerHTML = '<div class="rightbar-empty">Failed to load</div>';
        }
    }

    matchRoute(pattern, path) {
        const patternParts = pattern.split('/');
        const pathParts = path.split('/');

        if (patternParts.length !== pathParts.length) {
            return null;
        }

        const params = {};
        for (let i = 0; i < patternParts.length; i++) {
            if (patternParts[i].startsWith(':')) {
                params[patternParts[i].slice(1)] = pathParts[i];
            } else if (patternParts[i] !== pathParts[i]) {
                return null;
            }
        }

        return { params };
    }

    start() {
        this.resolve();
    }
}

export const router = new Router();