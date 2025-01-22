from flask import Blueprint, render_template, request, redirect, url_for, session
import requests

home_bp = Blueprint('home', __name__)

API_BASE_URL = 'http://localhost:8080/api'

@home_bp.route('/', methods=['GET'])
def redirect_home():
    return redirect(url_for('home.index'))

@home_bp.route('/home', methods=['GET'])
def index():
    logged_in = 'user_id' in session
    tag_ids = request.args.get('tags')
    user_id = session.get('user_id') if logged_in else None


    # Fetch all posts and tags from the API
    posts_response = requests.get(f"{API_BASE_URL}/posts/")
    posts_response.raise_for_status()
    posts = posts_response.json()

    tags_response = requests.get(f"{API_BASE_URL}/posts/tags/")
    tags_response.raise_for_status()
    tags = tags_response.json().get('tags', [])
    if not tags:
        tags = []

    # Fetch followers from the API
    followers_response = requests.get(f"{API_BASE_URL}/follow/followers/")
    followers_response.raise_for_status()
    followers_data = followers_response.json()

    # Extract the list of followers from the response
    followers = followers_data.get('followers', [])

    # Filter posts by tags if provided
    if tag_ids:
        # Split the comma-separated tag IDs and convert to a set for comparison
        tag_ids = set(tag_ids.split(','))

        # Filter posts that contain any of the selected tags
        filtered_posts = []
        if posts is None:
            posts = []
        for post in posts:
            post_tag_ids = {str(tag['id']) for tag in post.get('tags', [])}
            if post_tag_ids.intersection(tag_ids):
                filtered_posts.append(post)

        posts = filtered_posts  # Update the posts list with filtered posts

    # Create the follower map
    follower_map = {}
    if not posts:
        posts = []
    print(posts)
    for post in posts:
        # Sort the contents of the post
        post_id = post.get('id')
        post_contents = sorted(post.get('contents', []), key=lambda c: c.get('paragraph', 0))
        post['contents'] = post_contents  # Update the post with sorted contents

        likes_response = requests.get(f"{API_BASE_URL}/posts/{post_id}/likes")
        likes_response.raise_for_status()
        likes = likes_response.json().get('likes', [])
        post['count_likes'] = len(likes)

        # Fetch comments and their authors (Author details are already included in the comment response)
        comments_response = requests.get(f"{API_BASE_URL}/posts/{post_id}/comments")
        comments_response.raise_for_status()
        comments = comments_response.json().get('comments', [])

        # Count comments
        count_comments = len(comments)
        post['count_comments'] = count_comments


        is_liked = False
        if logged_in:
            like_check_response = requests.get(f"{API_BASE_URL}/posts/{post_id}/like/{user_id}")
            like_check_response.raise_for_status()
            is_liked = like_check_response.json().get("liked", False)
            post['is_liked'] = is_liked


        # Get the author ID for the current post
        author_id = post.get('author', {}).get('ID')
        response = requests.get(f"{API_BASE_URL}/follow/users/{author_id}/followers")
        if author_id:
            followers = response.json().get('followers', [])  # Assuming the API returns a list of user IDs
            if followers == None:
                post['is_following'] = False
            else:
                followers_users_id = [follower.get('id') for follower in followers]
                if user_id in followers_users_id:
                    post['is_following'] = True
                else:
                    post['is_following'] = False
            # Check if the current user has liked the post

    return render_template('index.html', posts=posts, tags=tags, logged_in=logged_in)
