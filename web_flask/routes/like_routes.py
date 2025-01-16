from flask import Blueprint, session, request, jsonify
import requests

like_bp = Blueprint('like', __name__, url_prefix='/like')

API_BASE_URL = "http://127.0.0.1:8080/api/posts"

@like_bp.route('/post/<post_id>', methods=['POST'])
def like_post(post_id):
    user_id = session.get('user_id')
    if not user_id:
        return {'success': False, 'message': 'Please log in to like posts.'}, 401

    # API endpoints
    like_api_url = f"{API_BASE_URL}/{post_id}/like"
    check_api_url = f"{API_BASE_URL}/{post_id}/like/{user_id}"
    likes_url = f"{API_BASE_URL}/{post_id}/likes"
    token = session.get('jwt_token')

    headers = {'Authorization': f'Bearer {token}'}
    likes = requests.get(likes_url, headers=headers).json().get('likes')

    # Step 1: Check if the user already liked the post
    check_response = requests.get(check_api_url, headers=headers)

    if check_response.status_code == 200:
        data = check_response.json()
        if data.get('liked'):  # User has already liked the post
            # Step 2: Unlike the post
            for like in likes:
                if like.get('author_id') == user_id:
                    like_id = like.get('id')
                    break
            unlike_response = requests.delete(f"{like_api_url}/{like_id}", headers=headers)
            if unlike_response.status_code == 200:
                likes = requests.get(likes_url, headers=headers).json().get('likes')
                like_count = len(likes)

                return {
                    'success': True,
                    'like_count': like_count,
                    'is_liked': False
                }, 200
            else:
                return {'success': False, 'message': 'Failed to unlike the post.'}, unlike_response.status_code
        else:
            # Step 3: Like the post
            like_response = requests.post(like_api_url, headers=headers)
            if like_response.status_code == 200:
                likes = requests.get(likes_url, headers=headers).json().get('likes')
                like_count = len(likes)

                return {
                    'success': True,
                    'like_count': like_count,
                    'is_liked': True
                }, 200
            else:
                return {'success': False, 'message': 'Failed to like the post.'}, like_response.status_code
    else:
        return {'success': False, 'message': 'Failed to fetch like status.'}, check_response.status_code
