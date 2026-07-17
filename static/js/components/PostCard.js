/**
 * PostCard component — displays a single post in the feed
 */
import { api } from '../api.js';
import { router } from '../router.js';
import { createAvatar } from './Avatar.js';
import { escapeHtml, timeAgo } from '../utils.js';

export function createPostCard(post) {
    const card = document.createElement('div');
    card.className = 'post-card';
    
    const avatar = createAvatar(post.nickname, 'sm');
    
    card.innerHTML = `
        <div class="post-header">
            <div class="post-user"></div>
            <button class="post-menu">...</button>
        </div>
        <div class="post-content">${escapeHtml(post.content)}</div>
        <div class="post-actions">
            <button class="action-btn like-btn" data-post-id="${post.postId}" data-score="1">
                <svg width="20" height="18" viewBox="0 0 20 18" fill="none">
                    <path d="M10 17C10 17 2 11.5 2 6C2 3.5 4 2 6 2C8 2 10 4 10 4C10 4 12 2 14 2C16 2 18 3.5 18 6C18 11.5 10 17 10 17Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span>${post.score || 0}</span>
            </button>
            <button class="action-btn comment-btn" data-post-id="${post.postId}">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                    <path d="M2 4C2 3.44772 2.44772 3 3 3H17C17.5523 3 18 3.44772 18 4V14C18 14.5523 17.5523 15 17 15H6L3 18V4Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span>${post.commentsCounter || 0}</span>
            </button>
        </div>
    `;
    
    // Insert avatar + user info into post-user
    const postUser = card.querySelector('.post-user');
    postUser.appendChild(avatar);
    
    const userInfo = document.createElement('div');
    userInfo.className = 'post-user-info';
    userInfo.innerHTML = `
        <span class="post-username">${escapeHtml((post.nickname || '').toUpperCase())}</span>
        <span class="post-time">${timeAgo(post.createdAt)}</span>
    `;
    postUser.appendChild(userInfo);
    
    // Click to view post detail
    card.addEventListener('click', (e) => {
        if (!e.target.closest('.action-btn') && !e.target.closest('.post-menu')) {
            router.navigate(`post/${post.postId}`);
        }
    });
    
    // Like button
    const likeBtn = card.querySelector('.like-btn');
    likeBtn.addEventListener('click', async (e) => {
        e.stopPropagation();
        try {
            const res = await api.react('post', post.postId, 1);
            likeBtn.querySelector('span').textContent = res.data.totalScore;
        } catch (err) {
            console.error('Failed to react:', err);
        }
    });
    
    // Comment button
    const commentBtn = card.querySelector('.comment-btn');
    commentBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        router.navigate(`post/${post.postId}`);
    });
    
    return card;
}