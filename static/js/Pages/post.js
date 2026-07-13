import { api } from '../api.js';

export async function renderPost(app, params) {

    
    const postDTO = await (async () => {
        try {
            return await api.getPost(parseInt(params.id));
        } catch (err) {
            app.innerHTML = err.message;
        }
    })();
    console.log(postDTO)
    
    if (postDTO.success) {
        app.innerHTML = `
            <div>
                <p>Title: ${postDTO.data.title}</p>
                <p>Body: ${postDTO.data.content}</p>
                <span>score: ${postDTO.data.score} - comments: ${postDTO.data.commentsCounter}</span>
            </div>            
        `;
    } else {
        app.innerHTML = `
            <p style="color: 'red'">${app.message}</p>
        `;
    }

}

