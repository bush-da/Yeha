from flask import Blueprint, request, redirect, url_for, flash, session
import requests

comment_bp = Blueprint('comment', __name__, url_prefix='/comments')

# Base URL for the API
API_BASE_URL = "http://127.0.0.1:8080/api/posts"

@comment_bp.route('/add/<post_id>', methods=['POST'])
def add_comment(post_id):
    """Add a comment to a post."""
    token = session.get('jwt_token')

    headers = {'Authorization': f'Bearer {token}'}

    user_id = session.get('user_id')
    if not user_id:
        flash('Please log in to comment.')
        return redirect(url_for('auth.login'))

    comment_text = request.form.get('comment')
    if comment_text:
        payload = {
            "author_id": user_id,
            "post_id": post_id,
            "content": comment_text
        }
        try:
            response = requests.post(f"{API_BASE_URL}/{post_id}/comment", json=payload, headers=headers)
            if response.status_code == 201:
                flash('Comment added successfully.')
            else:
                flash('Failed to add comment. Please try again.')
        except requests.RequestException as e:
            flash(f"An error occurred: {e}")
    else:
        flash('Comment cannot be empty.')

    return redirect(url_for('post.post_detail', post_id=post_id))

@comment_bp.route('/delete/<comment_id>/<post_id>', methods=['POST'])
def delete_comment(comment_id, post_id):
    """Delete a comment using the API."""

    token = session.get('jwt_token')

    headers = {'Authorization': f'Bearer {token}'}

    user_id = session.get('user_id')
    if not user_id:
        flash('Please log in to delete comments.')
        return redirect(url_for('auth.login'))

    try:
        # Send DELETE request to the API
        response = requests.delete(f"{API_BASE_URL}/{post_id}/comment/{comment_id}", headers=headers)
        if response.status_code == 204:
            flash('Comment deleted successfully.')
        else:
            flash('Failed to delete comment. Please try again.')
    except requests.RequestException as e:
        flash(f"An error occurred: {e}")

    # Redirect to the associated post page
    if post_id:
        return redirect(url_for('post.post_detail', post_id=post_id))
    else:
        flash('Post ID not found. Returning to the home page.')
        return redirect(url_for('home.index'))
