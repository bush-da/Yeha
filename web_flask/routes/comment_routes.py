from flask import Blueprint, request, redirect, url_for, flash, session
from models.comment import Comment
from models import storage

comment_bp = Blueprint('comment', __name__, url_prefix='/comments')

@comment_bp.route('/add/<post_id>', methods=['POST'])
def add_comment(post_id):
    """Add a comment to a post."""
    user_id = session.get('user_id')
    if not user_id:
        flash('Please log in to comment.')
        return redirect(url_for('auth.login'))

    comment_text = request.form.get('comment')
    if comment_text:
        new_comment = Comment(user_id=user_id, post_id=post_id, text=comment_text)
        storage.new(new_comment)
        storage.save()

        flash('Comment added successfully.')
    else:
        flash('Comment cannot be empty.')

    return redirect(url_for('post.post_detail', post_id=post_id))
