from flask import Blueprint, session, redirect, url_for, flash
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

    all_likes = storage.all(Like)
    like_key = f"Like.{user_id}_{post_id}"
    if like_key in all_likes:
        flash('You have already liked this post.')
        return redirect(request.referrer)

    new_like = Like(user_id=user_id, post_id=post_id)
    storage.new(new_like)
    storage.save()

    flash('Post liked successfully.')
    return redirect(request.referrer)
