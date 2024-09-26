from flask import Blueprint, session, redirect, url_for, flash, request
from models import storage
from models.user import User
from models.follower import Follower

follow_bp = Blueprint('follow', __name__)

@follow_bp.route('/follow/<user_id>', methods=['POST'])
def follow_user(user_id):
    """Follow or unfollow a user."""
    current_user_id = session.get('user_id')

    if not current_user_id:
        flash('Please log in to follow users.')
        return redirect(url_for('auth.login'))

    # Check if the current user is already following the target user
    existing_follow = None
    follows = storage.all(Follower).values()

    for follow in follows:
        if follow.follower_id == current_user_id and follow.followed_id == user_id:
            existing_follow = follow
            break

    # If already following, unfollow
    if existing_follow:
        storage.delete(existing_follow)
        storage.save()
        flash('User unfollowed successfully.')
    else:
        # Otherwise, follow the user
        new_follow = Follower(follower_id=current_user_id, followed_id=user_id)
        storage.new(new_follow)
        storage.save()
        flash('User followed successfully.')

    return redirect(request.referrer)
