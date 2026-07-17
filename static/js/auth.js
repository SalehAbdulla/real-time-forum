import { api } from './api.js';
import { router } from './router.js';

let currentTab = 'login';

export function renderLoginPage(app) {
    app.innerHTML = `
        <div class="auth-container">
            <div class="auth-logo">
                <span>forum</span><span class="dot"></span>
            </div>
            <div class="auth-tabs">
                <button class="auth-tab ${currentTab === 'login' ? 'active' : ''}" data-tab="login">Sign in</button>
                <button class="auth-tab ${currentTab === 'register' ? 'active' : ''}" data-tab="register">Register</button>
            </div>
            <div class="auth-title">${currentTab === 'login' ? 'Sign in to your account' : 'Create your account'}</div>
            <div id="auth-error" class="error-message" style="display: none;"></div>
            <form id="auth-form" class="auth-form">
                ${currentTab === 'login' ? getLoginFields() : getRegisterFields()}
            </form>
            <button type="submit" form="auth-form" class="btn-submit">${currentTab === 'login' ? 'Sign in' : 'Register'}</button>
            <div class="auth-footer">
                <span>${currentTab === 'login' ? "Don't have an account?" : 'Already have an account?'}</span>
                <a id="auth-switch">${currentTab === 'login' ? 'Register' : 'Sign in'}</a>
            </div>
            ${currentTab === 'login' ? `
            <div class="remember-me">
                <input type="checkbox" id="remember-me">
                <label for="remember-me">Remember Me</label>
            </div>` : ''}
        </div>
    `;

    document.querySelectorAll('.auth-tab').forEach(tab => {
        tab.addEventListener('click', () => {
            currentTab = tab.dataset.tab;
            renderLoginPage(app);
        });
    });

    document.getElementById('auth-switch').addEventListener('click', (e) => {
        e.preventDefault();
        currentTab = currentTab === 'login' ? 'register' : 'login';
        renderLoginPage(app);
    });

    document.getElementById('auth-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const errorEl = document.getElementById('auth-error');
        errorEl.style.display = 'none';

        const btn = document.querySelector('.btn-submit');
        btn.disabled = true;
        btn.textContent = 'Loading...';

        try {
            if (currentTab === 'login') {
                await handleLogin();
            } else {
                await handleRegister();
            }
        } catch (err) {
            errorEl.textContent = err.error || err.message || 'An error occurred';
            errorEl.style.display = 'block';
        } finally {
            btn.disabled = false;
            btn.textContent = currentTab === 'login' ? 'Sign in' : 'Register';
        }
    });
}

function getLoginFields() {
    return `
        <div class="form-group">
            <input type="text" name="identifier" class="form-input" placeholder="Email / Username" required>
        </div>
        <div class="form-group">
            <input type="password" name="password" class="form-input" placeholder="Password" required>
        </div>
    `;
}

function getRegisterFields() {
    return `
        <div class="form-group">
            <input type="text" name="nickName" class="form-input" placeholder="Nickname" required>
        </div>
        <div class="form-group">
            <input type="email" name="email" class="form-input" placeholder="Email" required>
        </div>
        <div class="form-group">
            <input type="text" name="firstName" class="form-input" placeholder="First Name" required>
        </div>
        <div class="form-group">
            <input type="text" name="lastName" class="form-input" placeholder="Last Name" required>
        </div>
        <div class="form-group">
            <input type="password" name="password" class="form-input" placeholder="Password" required>
        </div>
        <div class="form-group">
            <input type="password" name="confirmPassword" class="form-input" placeholder="Confirm Password" required>
        </div>
        <div class="form-group">
            <input type="text" name="age" class="form-input" placeholder="Age" required>
        </div>
        <div class="form-group">
            <input type="text" name="gender" class="form-input" placeholder="Gender" required>
        </div>
    `;
}

async function handleLogin() {
    const form = document.getElementById('auth-form');
    const identifier = form.querySelector('[name="identifier"]').value;
    const password = form.querySelector('[name="password"]').value;
    const rememberMe = document.getElementById('remember-me')?.checked || false;

    const res = await api.login(identifier, password, rememberMe);
    window.__isAuthenticated = true;
    // Fetch user profile
    try {
        const me = await api.me();
        window.__user = me.data;
    } catch {}
    router.navigate('feed');
}

async function handleRegister() {
    const form = document.getElementById('auth-form');
    const fields = {
        nickName: form.querySelector('[name="nickName"]').value,
        email: form.querySelector('[name="email"]').value,
        firstName: form.querySelector('[name="firstName"]').value,
        lastName: form.querySelector('[name="lastName"]').value,
        password: form.querySelector('[name="password"]').value,
        confirmPassword: form.querySelector('[name="confirmPassword"]').value,
        age: form.querySelector('[name="age"]').value,
        gender: form.querySelector('[name="gender"]').value,
    };

    const res = await api.register(fields);
    window.__isAuthenticated = true;
    // Fetch user profile
    try {
        const me = await api.me();
        window.__user = me.data;
    } catch {}
    router.navigate('feed');
}
