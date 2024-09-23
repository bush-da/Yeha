from flask import Blueprint, session, redirect, url_for, flash, request
from models import storage
from models.like import Like
from models.post import Post

like_bp = Blueprint('like', __name__, url_prefix='/like')

@like_bp.route('/post/<post_id>', methods=['GET', 'POST'])
def like_post(post_id):
    """Like a post."""
    user_id = session.get('user_id')
    if not user_id:
        flash('Please log in to like posts.')
        return redirect(url_for('auth.login'))

    all_likes = storage.all(Like).values()
    for likes in all_likes:
        if likes.author_id == user_id and likes.post_id == post_id:
            storage.delete(likes)
            storage.save()
            return redirect(request.referrer)


    new_like = Like(author_id=user_id, post_id=post_id)
    storage.new(new_like)
    storage.save()

    flash('Post liked successfully.')
    return redirect(request.referrer)
