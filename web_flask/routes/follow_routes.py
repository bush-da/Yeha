from flask import Blueprint, session, redirect, url_for, flash
from models import storage
from models.user import User
from models.follower import Follower

follow_bp = Blueprint('follow', __name__, url_prefix='/follow')

@follow_bp.route('/<int:user_id>')
def follow_user(user_id):
    """Follow a user."""
    current_user_id = session.get('user_id')
    if not current_user_id:
        flash('Please log in to follow users.')
        return redirect(url_for('auth.login'))

    # Check if already following
    follows = storage.all(Follower)
    for follow in follows.values():
        if follow.follower_id == current_user_id and follow.followed_id == user_id:
            flash('Already following this user.')
            return redirect(request.referrer)

    # Create a new follow record
    new_follow = Follower(follower_id=current_user_id, followed_id=user_id)
    storage.new(new_follow)
    storage.save()

    flash('User followed successfully.')
    return redirect(request.referrer)

@follow_bp.route('/unfollow/<int:user_id>')
def unfollow_user(user_id):
    """Unfollow a user."""
    current_user_id = session.get('user_id')
    if not current_user_id:
        flash('Please log in to unfollow users.')
        return redirect(url_for('auth.login'))

    follows = storage.all(Follower)
    for follow in follows.values():
        if follow.follower_id == current_user_id and follow.followed_id == user_id:
            storage.delete(follow)
            storage.save()
            flash('User unfollowed successfully.')
            return redirect(request.referrer)

    flash('You are not following this user.')
    return redirect(request.referrer)
