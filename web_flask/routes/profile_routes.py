from flask import Blueprint, render_template, redirect, url_for, request, flash, session
from werkzeug.utils import secure_filename
from werkzeug.security import generate_password_hash, check_password_hash
import os
import requests

ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg', 'gif'}

profile_bp = Blueprint('profile', __name__, url_prefix='/profile')

API_BASE_URL = 'http://localhost:8080/api'

@profile_bp.route('/<user_id>', methods=['GET'])
def profile(user_id):

    current_user_id = session.get('user_id')

    if not current_user_id:
        flash('Please log in to follow users.')
        return redirect(url_for('auth.login'))

    token = session.get('jwt_token')

    # Check if the post exists
    headers = {'Authorization': f'Bearer {token}'}

    user_response = requests.get(f"{API_BASE_URL}/users/{user_id}", headers=headers)
    user_response.raise_for_status()
    user = user_response.json().get('user')
    if not user:
        flash('User not found!')
        return redirect(url_for('home.index'))

    # Get user's posts
    posts_response = requests.get(f"{API_BASE_URL}/posts")
    posts_response.raise_for_status()
    posts = posts_response.json()

    if not posts:
        posts = []


    user_posts = [post for post in posts if post.get('author').get('ID') == user.get('ID') ]


    # Check if the user is viewing their own profile
    is_own_profile = session.get('user_id') == user_id
    # Fetch followers from the API

    response = requests.get(f"{API_BASE_URL}/follow/users/{user_id}/followers")
    followers = response.json().get('followers', [])  # Assuming the API returns a list of user IDs
    if followers == None:
        is_following = False
    else:
        followers_users_id = [follower.get('id') for follower in followers]
        if session.get('user_id') in followers_users_id:
            is_following = True
        else:
            is_following = False

    followers_response = requests.get(f"{API_BASE_URL}/follow/users/{user_id}/followers/")
    followers_response.raise_for_status()
    if followers_response.json().get('followers') == None:
        follower_count = 0
    else:
        follower_count = len(followers_response.json().get('followers'))
    following_response = requests.get(f"{API_BASE_URL}/follow/users/{user_id}/following/")
    following_response.raise_for_status()
    if following_response.json().get('following') == None:
        following_count = 0
    else:
        following_count = len(following_response.json().get('following'))

    return render_template('profile.html', user=user, posts=user_posts, is_own_profile=is_own_profile,follower_count=follower_count, following_count=following_count, is_following=is_following)

@profile_bp.route('/<user_id>/update', methods=['POST'])
def update_profile(user_id):
    token = session.get('jwt_token')  # Get the JWT token

    headers = {'Authorization': f'Bearer {token}'}
    url = f"{API_BASE_URL}/users/{user_id}"

    # Check if the user exists and validate ownership
    user_response = requests.get(url, headers=headers)
    if user_response.status_code != 200:
        flash('User not found or unauthorized access!')
        return redirect(url_for('profile.profile', user_id=user_id))

    # Extract updated values from form data
    payload = {
        "username": request.form.get('username'),
        "email": request.form.get('email'),
        "bio": request.form.get('bio'),
        "gender": request.form.get('gender'),
    }

    # Handle profile picture upload
    profile_picture = request.files.get('profile_picture')

    if profile_picture:
        filename = save_profile_picture(profile_picture)
        payload["profile_picture"] = filename

    # Call the API to update the profile
    response = requests.put(url, json=payload, headers=headers)

    if response.status_code == 200:
        flash('Profile updated successfully!')
    else:
        flash('Failed to update profile!')

    return redirect(url_for('profile.profile', user_id=user_id))


@profile_bp.route('/<user_id>/update-password', methods=['GET'])
def update_password_form(user_id):
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}
    url = f"{API_BASE_URL}/users/{user_id}"

    # Check if the user exists
    user_response = requests.get(url, headers=headers)
    if user_response.status_code != 200:
        flash('User not found or unauthorized access!')
        return redirect(url_for('profile.profile', user_id=user_id))

    return render_template('update_password.html', user_id=user_id)


@profile_bp.route('/<user_id>/update-password', methods=['POST'])
def update_password(user_id):
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}
    url = f"{API_BASE_URL}/users/{user_id}/password"

    # Extract password fields from the form
    old_password = request.form.get('old_password')
    new_password = request.form.get('new_password')
    confirm_password = request.form.get('confirm_password')

    # Validation
    if not old_password or not new_password or not confirm_password:
        flash("All fields are required!")
        return redirect(url_for('profile.update_password', user_id=user_id))

    if new_password != confirm_password:
        flash("New passwords do not match!")
        return redirect(url_for('profile.update_password', user_id=user_id))

    # Call API to update password
    payload = {
        "old_password": old_password,
        "new_password": new_password
    }
    response = requests.put(url, json=payload, headers=headers)

    if response.status_code == 200:
        flash("Password updated successfully!")
    else:
        flash("Failed to update password. " + response.json().get("error", "Unknown error."))

    return redirect(url_for('profile.profile', user_id=user_id))

def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS


def save_profile_picture(profile_picture):
    if not allowed_file(profile_picture.filename):
        raise ValueError("Invalid file type")

    filename = secure_filename(profile_picture.filename)
    filepath = os.path.join('/home/smuca/projects/Yeha/web_flask/static/images', filename)

    try:
        profile_picture.save(filepath)
        if os.path.exists(filepath):
            print(f"File successfully saved at: {filepath}")
        else:
            print(f"File failed to save at: {filepath}")
    except Exception as e:
        print(f"Error saving file: {e}")
        raise

    return filename
