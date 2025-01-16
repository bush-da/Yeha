from flask import Blueprint, render_template, request, redirect, url_for, flash, session
import requests
from werkzeug.utils import secure_filename
import os

post_bp = Blueprint('post', __name__)

API_BASE_URL = 'http://localhost:8080/api'

@post_bp.route('/posts/new_post', methods=['GET', 'POST'])
def new_post():
    if not session.get('user_id'):
        return redirect(url_for('home.index'))

    if request.method == 'POST':
        title = request.form.get('title')
        author_id = session.get('user_id')
        tags_input = request.form.get('tags').split(',')
        category_id = request.form.get('category')

        if not title:
            flash('Title is required!')
            return redirect(url_for('post.new_post'))

        # Prepare data for API request
        contents = []
        block_index = 0
        while True:
            block_type = request.form.get(f'block_type_{block_index}')
            if not block_type:
                break

            if block_type == 'text':
                text_content = request.form.get(f'text_content_{block_index}')
                if text_content:
                    contents.append({
                        "content": text_content,
                        "paragraph": block_index
                    })

            elif block_type == 'image':
                image = request.files.get(f'image_content_{block_index}')
                if image and image.filename:
                    filename = secure_filename(image.filename)
                    file_path = os.path.join('/home/smuca/projects/Yeha/web_flask/static/images', filename)

                    # Save the image file
                    image.save(file_path)

                    # Generate the URL for the saved image
                    file_url = url_for('static', filename='images/' + filename)

                    contents.append({
                        "file_url": file_url,
                        "paragraph": block_index
                    })

            block_index += 1

        # Prepare tags
        tags = [{"name": tag_name.strip().lower()} for tag_name in tags_input if tag_name.strip()]

        # API request payload
        payload = {
            "title": title,
            "contents": contents,
            "tags": tags
        }

        # Send API request

        headers = {"Authorization": f"Bearer {session.get('jwt_token')}"}  # Add auth token if required
        response = requests.post(f"{API_BASE_URL}/posts/posts", json=payload, headers=headers)

        if response.status_code == 201:
            flash('Post created successfully!')
            return redirect(url_for('home.index'))
        else:
            flash(f"Error: {response.json().get('error', 'Failed to create post')}")
            return redirect(url_for('post.new_post'))

    # Fetch categories for the form (if required for the UI)
    # This step assumes you have a Go API endpoint to fetch categories.
    # categories_response = requests.get(f"{GO_API_BASE_URL}/categories")
    # categories = categories_response.json() if categories_response.status_code == 200 else []

    return render_template('new_post.html')

@post_bp.route('/post/<post_id>/update', methods=['GET', 'POST'])
def update_post(post_id):
    try:
        # Fetch the post details
        post_response = requests.get(f"{API_BASE_URL}/posts/posts/{post_id}")
        post_response.raise_for_status()
        post = post_response.json().get("post")
        headers = {"Authorization": f"Bearer {session.get('jwt_token')}"}  # Add auth token if required

        if not post:
            raise ValueError("Invalid API response: Missing 'post' data")

        if request.method == 'POST':
            # Gather title and tags
            title = request.form.get('title')
            tags_input = request.form.get('tags').split(',')
            tags = [{"name": tag.strip().lower()} for tag in tags_input if tag.strip()]

            # Send updated title and tags in a single request
            payload = {"title": title, "tags": tags}
            requests.put(f"{API_BASE_URL}/posts/posts/{post_id}", json=payload, headers=headers)

            # Update content blocks
            block_index = 0
            while True:
                content_id = request.form.get(f'content_id_{block_index}')
                block_type = request.form.get(f'block_type_{block_index}')
                if not block_type:
                    break

                if content_id and request.form.get(f'delete_content_{content_id}') == 'true':
                    requests.delete(f"{API_BASE_URL}/posts/posts/contents/{content_id}", headers=headers)
                else:
                    content_data = {"post_id": post_id, "paragraph": block_index}
                    if block_type == 'text':
                        content_data["content"] = request.form.get(f'text_content_{block_index}')
                    elif block_type == 'image':
                        image = request.files.get(f'image_content_{block_index}')
                        if image and image.filename:
                            filename = secure_filename(image.filename)
                            file_path = os.path.join('/home/smuca/projects/Yeha/web_flask/static/images', filename)
                            image.save(file_path)
                            content_data["file_url"] = url_for('static', filename='images/' + filename)

                    if content_id:
                        requests.put(f"{API_BASE_URL}/posts/posts/contents/{content_id}", json=content_data, headers=headers)
                    else:
                        requests.post(f"{API_BASE_URL}/posts/posts/contents/", json=content_data, headers=headers)

                block_index += 1

            flash('Post updated successfully!')
            return redirect(url_for('profile.profile', user_id=post.get("author").get("ID")))

        # Fetch post contents
        content_response = requests.get(f"{API_BASE_URL}/posts/contents/post/{post_id}")
        content_response.raise_for_status()
        contents = sorted(content_response.json(), key=lambda c: c.get("paragraph"))


        return render_template('update_post.html', post=post, contents=contents)
    except requests.RequestException as e:
        flash(f"Error updating post: {str(e)}", 'error')
        return redirect(url_for('home.index'))

# Delete Post
@post_bp.route('/post/<post_id>/delete', methods=['POST'])
def delete_post(post_id):
    try:
        # Get the current user's JWT token (stored in session or elsewhere)
        token = session.get('jwt_token')

        # Check if the post exists
        headers = {'Authorization': f'Bearer {token}'}
        post_response = requests.get(f"{API_BASE_URL}/posts/posts/{post_id}", headers=headers)
        post_response.raise_for_status()
        post = post_response.json()

        # Delete the post
        delete_response = requests.delete(f"{API_BASE_URL}/posts/posts/{post_id}", headers=headers)
        delete_response.raise_for_status()

        flash('Post deleted successfully!')
        post = post.get("post")
        user_id = post.get("author").get('ID')
        return redirect(url_for('profile.profile', user_id=user_id))
    except requests.RequestException as e:
        flash(f"Error deleting post: {str(e)}", 'error')
        return redirect(url_for('home.index'))


@post_bp.route('/posts/<post_id>', methods=['GET'])
def post_detail(post_id):
    try:
        logged_in = 'user_id' in session
        user_id = session.get('user_id') if logged_in else None

        # Fetch post details from the API
        post_response = requests.get(f"{API_BASE_URL}/posts/posts/{post_id}")
        post_response.raise_for_status()
        post_data = post_response.json()

        # Validate response structure
        post = post_data.get("post")
        if not post:
            raise ValueError("Invalid API response: Missing 'post' data")

        count_likes = post.get('likes')

        # Check if the current user has liked the post
        is_liked = False
        if logged_in:
            like_check_response = requests.get(f"{API_BASE_URL}/posts/{post_id}/like/{user_id}")
            like_check_response.raise_for_status()
            is_liked = like_check_response.json().get("liked", False)

        # Fetch comments and their authors (Author details are already included in the comment response)
        comments_response = requests.get(f"{API_BASE_URL}/posts/{post_id}/comments")
        comments_response.raise_for_status()
        comments = comments_response.json().get('comments', [])

        # Count comments
        count_comments = len(comments)
        author_id = post.get("author").get('ID')

        # Map followers for the post's author
        follower_map = list()
        followers_response = requests.get(f"{API_BASE_URL}/follow/users/{author_id}/followers")
        followers = followers_response.json().get('followers', [])

        # Iterate through the list of followers and extract 'follower_id'
        for follower in followers:
            follower_map.append(follower.get('follower_id'))


        post['contents'] = sorted(post['contents'], key=lambda c: c.get("paragraph"))


        return render_template(
            'post_detail.html',
            post=post,
            count_likes=count_likes,
            count_comments=count_comments,
            comments=comments,
            follower_map=follower_map,
            logged_in=logged_in,
            is_liked=is_liked  # Pass the liked status to the template
        )
    except requests.RequestException as e:
        flash(f"Error fetching post details: {str(e)}", 'error')
        return redirect(url_for('home.index'))
    except ValueError as e:
        flash(str(e), 'error')
        return redirect(url_for('home.index'))
