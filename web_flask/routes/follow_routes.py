from flask import Blueprint, session, redirect, url_for, flash, request, render_template
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


@follow_bp.route('/<user_id>/followers', methods=['GET'])
def followers_list(user_id):
    """List followers of a user."""
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found!')
        return redirect(url_for('home.index'))

    # Get all followers
    followers = [f for f in storage.all(Follower).values() if f.followed_id == user_id]
    follower_users = [storage.all(User).get(f"User.{f.follower_id}") for f in followers]

    return render_template('followers_list.html', user=user, followers=follower_users)


@follow_bp.route('/<user_id>/following', methods=['GET'])
def following_list(user_id):
    """List users a user is following."""
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found!')
        return redirect(url_for('home.index'))

    # Get all following
    following = [f for f in storage.all(Follower).values() if f.follower_id == user_id]
    following_users = [storage.all(User).get(f"User.{f.followed_id}") for f in following]

    return render_template('following_list.html', user=user, following=following_users)
