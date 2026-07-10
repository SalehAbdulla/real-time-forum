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

        // Match route
        for (const pattern in this.routes) {
            const match = this.matchRoute(pattern, path);
            if (match) {
                // Check guard
                const guard = this.guards[pattern];
                if (guard && !guard()) {
                    this.navigate('/login');
                    return;
                }

                this.currentRoute = pattern;
                const app = document.getElementById('app');
                app.innerHTML = '';
                this.routes[pattern](app, match.params, queryString);
                return;
            }
        }

        // Fallback to login
        this.navigate('/login');
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
