from flask import Blueprint, render_template, request, redirect, url_for, flash, session
from models.post import Post
from models.content import Content
from models.category import Category
from models.tag import Tag
from models.like import Like
from models.comment import Comment
from models.follower import Follower
from models import storage
from werkzeug.utils import secure_filename
import os

post_bp = Blueprint('post', __name__)

@post_bp.route('/posts/', methods=['GET'])
def list_posts():
    posts = storage.all(Post).values()
    # Fetch contents for each post
    for post in posts:
        post.contents = list(storage.all(Content).values())
        post.contents = [content for content in post.contents if content.post_id == post.id]
    return render_template('index.html', posts=posts)

@post_bp.route('/posts/new_post', methods=['GET', 'POST'])
def new_post():
    if not session.get('user_id'):
        return redirect(url_for('home.index'))

    if request.method == 'POST':
        title = request.form.get('title')
        author_id = session.get('user_id')  # Assuming user ID is stored in session
        tags_input = request.form.get('tags').split(',')
        category_id = request.form.get('category')

        if not title:
            flash('Title is required!')
            return redirect(url_for('post.new_post'))

        # Create new post
        new_post = Post(title=title, author_id=author_id, category_id=category_id)
        storage.new(new_post)

        block_index = 0
        while True:
            block_type = request.form.get(f'block_type_{block_index}')
            if not block_type:
                break

            if block_type == 'text':
                text_content = request.form.get(f'text_content_{block_index}')
                if text_content:
                    # Store text content as paragraph(s)
                    new_content = Content(file_url=None, post_id=new_post.id, paragraph=block_index, content=text_content)
                    storage.new(new_content)

            elif block_type == 'image':
                image = request.files.get(f'image_content_{block_index}')
                if image and image.filename:
                    filename = secure_filename(image.filename)
                    file_path = os.path.join('/home/smuca/projects/Yeha/web_flask/static/images', filename)

                    # Save the image file
                    image.save(file_path)

                    # Generate the URL for the saved image
                    file_url = url_for('static', filename='images/' + filename)

                    new_content = Content(file_url=file_url, post_id=new_post.id, paragraph=block_index)
                    storage.new(new_content)

            block_index += 1

        # Handle tags
        for tag_name in tags_input:
            tag_name = tag_name.strip()
            if tag_name:
                tag = next((tag for tag in storage.all(Tag).values() if tag.name == tag_name.lower()), None)
                if not tag:
                    tag = Tag(name=tag_name.lower())
                    storage.new(tag)
                if tag not in new_post.tags:
                    new_post.tags.append(tag)

        storage.save()
        flash('Post created successfully!')
        return redirect(url_for('post.list_posts'))

    # Fetch categories for the form
    categories = storage.all(Category).values()
    return render_template('new_post.html', categories=categories)

# Update Post Route
@post_bp.route('/post/<post_id>/update', methods=['GET', 'POST'])
def update_post(post_id):
    post = storage.all(Post).get(f'Post.{post_id}')
    if request.method == 'POST':
        post.title = request.form.get('title')
        post.contents.content = request.form.get('contents')
        tags_input = request.form.get('tags').split(',')
        post.tags.clear()

        for tag_name in tags_input:
            tag_name = tag_name.strip()
            if tag_name:
                tag = next((tag for tag in storage.all(Tag).values() if tag.name == tag_name), None)
                if not tag:
                    tag = Tag(name=tag_name)
                    storage.new(tag)
                post.tags.append(tag)

        block_index = 0
        while True:
            content_id = request.form.get(f'content_id_{block_index}')
            block_type = request.form.get(f'block_type_{block_index}')

            if not block_type:
                break

            if content_id and request.form.get(f'delete_content_{content_id}') == 'true':
                content_to_delete = storage.all(Content).get(f'Content.{content_id}')
                if content_to_delete:
                    storage.delete(content_to_delete)
                block_index += 1
                continue

            if content_id:
                content = storage.all(Content).get(f'Content.{content_id}')
            else:
                content = None

            if block_type == 'text':
                text_content = request.form.get(f'text_content_{block_index}')
                if text_content:
                    if content:
                        content.content = text_content
                    else:
                        new_content = Content(file_url=None, post_id=post.id, paragraph=block_index, content=text_content)
                        storage.new(new_content)

            elif block_type == 'image':
                image = request.files.get(f'image_content_{block_index}')
                if image and image.filename:
                    filename = secure_filename(image.filename)
                    file_path = os.path.join('/home/smuca/projects/Yeha/web_flask/static/images', filename)
                    image.save(file_path)
                    file_url = url_for('static', filename='images/' + filename)

                    if content:
                        content.file_url = file_url
                    else:
                        new_content = Content(file_url=file_url, post_id=post.id, paragraph=block_index)
                        storage.new(new_content)

            block_index += 1

        storage.save()
        return redirect(url_for('profile.profile', user_id=post.author.id))
    post.contents = sorted(post.contents, key=lambda c: c.paragraph)
    return render_template('update_post.html', post=post)

# Delete Post Route
@post_bp.route('/post/<post_id>/delete', methods=['POST'])
def delete_post(post_id):
    post = storage.all(Post).get(f"Post.{post_id}")
    storage.delete(post)
    storage.save()
    return redirect(url_for('profile.profile', user_id=post.author.id))


@post_bp.route('/posts/<post_id>', methods=['GET'])
def post_detail(post_id):
    logged_in = 'user_id' in session
    post = storage.all(Post).get(f"Post.{post_id}")
    follower_map = {}
    if post:
        post.contents = sorted(post.contents, key=lambda c: c.paragraph)
        comments = [comment for comment in storage.all(Comment).values() if comment.post_id == post.id]
        count_likes = len([like for like in storage.all(Like).values() if like.post_id == post.id])
        followers = storage.all(Follower).values()
        # Filter followers for the current post's author
        follower_map[post.author.id] = [follower.follower_id for follower in followers if follower.followed_id == post.author.id]

        return render_template('post_detail.html', post=post, count_likes=count_likes, comments=comments, follower_map=follower_map, logged_in=logged_in)
    else:
        flash('Post not found', 'error')
        return redirect(url_for('home.index'))
