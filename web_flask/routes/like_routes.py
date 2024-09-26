from flask import Blueprint, session, redirect, url_for, flash, request
from models import storage
from models.like import Like
from models.post import Post

like_bp = Blueprint('like', __name__, url_prefix='/like')
@like_bp.route('/post/<post_id>', methods=['POST'])
def like_post(post_id):
    user_id = session.get('user_id')
    if not user_id:
        return {'success': False, 'message': 'Please log in to like posts.'}, 401

    # Check if the post exists
    post = storage.all(Post).get(f'Post.{post_id}')
    if not post:
        return {'success': False, 'message': 'Post not found.'}, 404

    # Query to check if the user already liked the post
    all_likes = storage.all(Like).values()
    existing_like = None

    for like in all_likes:
        if like.author_id == user_id and like.post_id == post_id:
            existing_like = like
            break

    if existing_like:
        # User has already liked, so we remove the like
        storage.delete(existing_like)
        storage.save()
        liked = False
    else:
        # User has not liked yet, so we add a new like
        new_like = Like(author_id=user_id, post_id=post_id)
        storage.new(new_like)
        storage.save()
        liked = True

    response_data = {
        'success': True,
        'like_count': post.count_likes(),
        'is_liked': liked
    }
    return response_data, 200  # Return the JSON response
