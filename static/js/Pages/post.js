import { api } from '../api.js';
import { escapeHtml, timeAgo } from '../utils.js';

export async function renderPost(app, params) {
    app.innerHTML = '<div class="loading-spinner">Loading post...</div>';
    
    try {
        const postDTO = await api.getPost(parseInt(params.id));
        
        if (!postDTO.success) {
            app.innerHTML = '<div class="error-state">Failed to load post</div>';
            return;
        }
        
        const post = postDTO.data;
        
        app.innerHTML = `
            <div class="post-detail">
                <div class="post-detail-header">
                    <button class="back-btn" id="back-btn">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                            <path d="M19 12H5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M12 19L5 12L12 5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                    </button>
                    <span class="post-detail-title">Post</span>
                </div>
                
                <div class="post-detail-card">
                    <div class="post-detail-user">
                        <div class="post-avatar">${getInitials(post.nickname)}</div>
                        <div class="post-detail-user-info">
                            <span class="post-username">${escapeHtml((post.nickname || '').toUpperCase())}</span>
                            <span class="post-time">${timeAgo(post.createdAt)}</span>
                        </div>
                    </div>
                    
                    <h2 class="post-detail-heading">${escapeHtml(post.title)}</h2>
                    <p class="post-detail-body">${escapeHtml(post.content)}</p>
                    
                    <div class="post-detail-actions">
                        <button class="action-btn like-btn" id="detail-like-btn">
                            <svg width="24" height="22" viewBox="0 0 20 18" fill="none">
                                <path d="M10 17C10 17 2 11.5 2 6C2 3.5 4 2 6 2C8 2 10 4 10 4C10 4 12 2 14 2C16 2 18 3.5 18 6C18 11.5 10 17 10 17Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                            <span>${post.score || 0}</span>
                        </button>
                        <button class="action-btn comment-btn" id="detail-comment-btn">
                            <svg width="24" height="24" viewBox="0 0 20 20" fill="none">
                                <path d="M2 4C2 3.44772 2.44772 3 3 3H17C17.5523 3 18 3.44772 18 4V14C18 14.5523 17.5523 15 17 15H6L3 18V4Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                            <span>${post.commentsCounter || 0}</span>
                        </button>
                    </div>
                </div>
                
                <div class="comments-section">
                    <h3 class="comments-heading">Comments</h3>
                    <div class="comment-form">
                        <textarea class="comment-input" id="comment-input" placeholder="Write a comment..." rows="3"></textarea>
                        <button class="comment-submit-btn" id="comment-submit-btn">Post Comment</button>
                    </div>
                    <div class="comments-list" id="comments-list">
                        <div class="loading-spinner">Loading comments...</div>
                    </div>
                </div>
            </div>
        `;
        
        // Back button
        document.getElementById('back-btn').addEventListener('click', () => {
            window.history.back();
        });
        
        // Like button - use userScore from API for initial state, then toggle
        const likeBtn = document.getElementById('detail-like-btn');
        // Set initial state from server data
        if (post.userScore === 1) {
            likeBtn.classList.add('liked');
        }
        likeBtn.addEventListener('click', async () => {
            try {
                const oldScore = parseInt(likeBtn.querySelector('span').textContent) || 0;
                const res = await api.react('post', post.postId, 1);
                const newScore = res.data.totalScore;
                likeBtn.querySelector('span').textContent = newScore;
                if (newScore < oldScore) {
                    likeBtn.classList.remove('liked');
                } else {
                    likeBtn.classList.add('liked');
                }
            } catch (err) {
                console.error('Failed to react:', err);
            }
        });
        
        // Comment submit
        const commentInput = document.getElementById('comment-input');
        const commentSubmitBtn = document.getElementById('comment-submit-btn');
        
        commentSubmitBtn.addEventListener('click', async () => {
            const content = commentInput.value.trim();
            if (!content || content.length < 1) return;
            
            commentSubmitBtn.disabled = true;
            commentSubmitBtn.textContent = 'Posting...';
            
            try {
                await api.createComment(post.postId, content);
                commentInput.value = '';
                // Reload comments
                loadComments(post.postId);
            } catch (err) {
                console.error('Failed to create comment:', err);
            } finally {
                commentSubmitBtn.disabled = false;
                commentSubmitBtn.textContent = 'Post Comment';
            }
        });
        
        // Load comments
        loadComments(post.postId);
        
    } catch (err) {
        app.innerHTML = `<div class="error-state">${escapeHtml(err.message || 'Failed to load post')}</div>`;
    }
}

async function loadComments(postId) {
    const commentsList = document.getElementById('comments-list');
    if (!commentsList) return;
    
    try {
        const res = await api.getComments(postId);
        const comments = res.data?.comments || [];
        
        if (comments.length === 0) {
            commentsList.innerHTML = '<div class="empty-state">No comments yet. Be the first to comment!</div>';
            return;
        }
        
        commentsList.innerHTML = comments.map(comment => {
            return `
            <div class="comment-card">
                <div class="comment-header">
                    <div class="comment-avatar">${getInitials(comment.nickname)}</div>
                    <div class="comment-info">
                        <span class="comment-username">${escapeHtml((comment.nickname || '').toUpperCase())}</span>
                        <span class="comment-time">${timeAgo(comment.createdAt)}</span>
                    </div>
                </div>
                <p class="comment-body">${escapeHtml(comment.commentText || '')}</p>
                <div class="comment-actions">
                    <button class="action-btn like-btn comment-like-btn" data-comment-id="${comment.commentId}">
                        <svg width="16" height="14" viewBox="0 0 20 18" fill="none">
                            <path d="M10 17C10 17 2 11.5 2 6C2 3.5 4 2 6 2C8 2 10 4 10 4C10 4 12 2 14 2C16 2 18 3.5 18 6C18 11.5 10 17 10 17Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                        <span>${comment.score || 0}</span>
                    </button>
                </div>
            </div>
        `}).join('');
        
        // Add like handlers for comments - derive state from score change
        document.querySelectorAll('.comment-like-btn').forEach(btn => {
            btn.addEventListener('click', async (e) => {
                e.stopPropagation();
                const commentId = parseInt(btn.dataset.commentId);
                try {
                    const oldScore = parseInt(btn.querySelector('span').textContent) || 0;
                    const res = await api.react('comment', commentId, 1);
                    const newScore = res.data.totalScore;
                    btn.querySelector('span').textContent = newScore;
                    if (newScore < oldScore) {
                        btn.classList.remove('liked');
                    } else {
                        btn.classList.add('liked');
                    }
                } catch (err) {
                    console.error('Failed to react on comment:', err);
                }
            });
        });
        
        // Update comment count in the header
        const commentCountSpan = document.querySelector('#detail-comment-btn span');
        if (commentCountSpan) {
            commentCountSpan.textContent = comments.length;
        }
        
    } catch (err) {
        commentsList.innerHTML = '<div class="empty-state">Failed to load comments.</div>';
    }
}

function getInitials(name) {
    if (!name) return '?';
    return name.substring(0, 2).toUpperCase();
}