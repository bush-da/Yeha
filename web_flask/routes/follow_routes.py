from flask import Blueprint, session, redirect, url_for, flash, request, render_template
import requests

follow_bp = Blueprint('follow', __name__)

# Base URL for the API
API_BASE_URL = "http://127.0.0.1:8080/api"


@follow_bp.route('/follow/<user_id>', methods=['POST'])
def follow_user(user_id):
    """Follow or unfollow a user using the API."""
    token = session.get('jwt_token')

    headers = {'Authorization': f'Bearer {token}'}

    current_user_id = session.get('user_id')

    if not current_user_id:
        flash('Please log in to follow users.')
        return redirect(url_for('auth.login'))

    # Check if the current user is already following the target user
    try:
        response = requests.get(f"{API_BASE_URL}/follow/users/{current_user_id}/following")
        if response.status_code == 200:
            following = response.json().get('following')  # Assuming the API returns a list of user IDs
            if following is None:
                following = []
            following_users_id = [followin.get('id') for followin in following]
            if user_id in following_users_id:
                # Unfollow the user
                unfollow_response = requests.delete(f"{API_BASE_URL}/follow/users/{user_id}/follow/", json={"follower_id": current_user_id}, headers=headers)
                if unfollow_response.status_code == 200:
                    return {
                        'success': True,
                        'is_following': False,
                    }, 200
                else:
                    flash('Failed to unfollow user. Please try again.')
            else:
                # Follow the user
                follow_response = requests.post(f"{API_BASE_URL}/follow/users/{user_id}/follow/", json={"follower_id": current_user_id}, headers=headers)
                if follow_response.status_code == 200:
                    return {
                        'success': True,
                        'is_following': True,
                    }, 200
                else:
                    flash('Failed to follow user. Please try again.')
        else:
            flash('Failed to fetch following list. Please try again.')
    except requests.RequestException as e:
        flash(f"An error occurred: {e}")

    return redirect(request.referrer)


@follow_bp.route('/<user_id>/followers', methods=['GET'])
def followers_list(user_id):
    """List followers of a user using the API."""
    token = session.get('jwt_token')

    headers = {'Authorization': f'Bearer {token}'}

    try:
        # Fetch user details (if needed for the template)
        user_response = requests.get(f"{API_BASE_URL}/users/{user_id}", headers=headers)
        if user_response.status_code != 200:
            flash('User not found!')
            return redirect(url_for('home.index'))
        user = user_response.json()

        # Fetch followers
        followers_response = requests.get(f"{API_BASE_URL}/follow/users/{user_id}/followers")
        if followers_response.status_code == 200:
            followers = followers_response.json().get('followers')  # Assuming API returns a list of follower user details
            if followers == None:
                followers = []
        else:
            flash('Failed to fetch followers. Please try again.')
            followers = []
    except requests.RequestException as e:
        flash(f"An error occurred: {e}")
        return redirect(url_for('home.index'))
    user = user.get('user')

    return render_template('followers_list.html', user=user, followers=followers)


@follow_bp.route('/<user_id>/following', methods=['GET'])
def following_list(user_id):
    """List users a user is following using the API."""
    token = session.get('jwt_token')

    headers = {'Authorization': f'Bearer {token}'}

    try:
        # Fetch user details (if needed for the template)
        user_response = requests.get(f"{API_BASE_URL}/users/{user_id}", headers=headers)
        if user_response.status_code != 200:
            flash('User not found!')
            return redirect(url_for('home.index'))
        user = user_response.json()

        # Fetch following
        following_response = requests.get(f"{API_BASE_URL}/follow/users/{user_id}/following")
        if following_response.status_code == 200:
            following = following_response.json().get("following", [])
            if following == None:
                following = []

        else:
            flash('Failed to fetch following list. Please try again.')
            following = []
    except requests.RequestException as e:
        flash(f"An error occurred: {e}")
        return redirect(url_for('home.index'))
    user = user.get('user')

    return render_template('following_list.html', user=user, following=following)
