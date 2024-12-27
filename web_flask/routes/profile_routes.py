from flask import Blueprint, render_template, redirect, url_for, request, flash, session
from werkzeug.utils import secure_filename
from models import storage
from models.user import User
from models.post import Post
from models.follower import Follower
from werkzeug.security import generate_password_hash, check_password_hash
from uuid import UUID
import os

profile_bp = Blueprint('profile', __name__, url_prefix='/profile')

@profile_bp.route('/<user_id>', methods=['GET'])
def profile(user_id):
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found!')
        return redirect(url_for('home.index'))

    # Get user's posts
    posts = storage.all(Post).values()
    user_posts = [post for post in posts if post.author_id == user.id]

    # Check if the user is viewing their own profile
    is_own_profile = session.get('user_id') == user_id
    followers = storage.all(Follower).values()
    follower_map = {}
    follower_map[user.id] = [follower.follower_id for follower in followers if follower.followed_id == user.id]
    # Followers and Following counts
    followers = [f for f in storage.all(Follower).values() if f.followed_id == user.id]
    following = [f for f in storage.all(Follower).values() if f.follower_id == user.id]

    follower_count = len(followers)
    following_count = len(following)

    return render_template('profile.html', user=user, posts=user_posts, is_own_profile=is_own_profile, follower_map=follower_map,follower_count=follower_count, following_count=following_count)

@profile_bp.route('/<user_id>/update', methods=['POST'])
def update_profile(user_id):
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found!')
        return redirect(url_for('profile.profile', user_id=user_id))

    if session.get('user_id') != user_id:
        flash('Unauthorized to update this profile!')
        return redirect(url_for('profile.profile', user_id=user_id))

    username = request.form.get('username')
    email = request.form.get('email')
    bio = request.form.get('bio')
    profile_picture = request.files.get('profile_picture')

    if username:
        user.username = username
    if email:
        user.email = email
    if bio:
        user.bio = bio
    if profile_picture:
        user.profile_picture = save_profile_picture(profile_picture)  # Assuming this function saves the picture and returns the filename

    storage.save()
    flash('Profile updated successfully!')
    return redirect(url_for('profile.profile', user_id=user_id))

def save_profile_picture(profile_picture):
    # Save the profile picture and return the filename
    # Ensure you handle file uploads and security properly
    filename = secure_filename(profile_picture.filename)
    filepath = os.path.join('/home/smuca/projects/Yeha/web_flask/static/images', filename)
    profile_picture.save(filepath)
    return filename
