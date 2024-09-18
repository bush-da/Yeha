from flask import Blueprint, render_template, request, redirect, url_for
from models.post import Post
from models.tag import Tag
from models import storage

home_bp = Blueprint('home', __name__)

@home_bp.route('/', methods=['GET'])
def redirect_home():
    return redirect(url_for('home.index'))

@home_bp.route('/home', methods=['GET'])
def index():
    tag_ids = request.args.get('tags')

    if tag_ids:
        # Convert the comma-separated string into a list of tag IDs
        tag_ids = tag_ids.split(',')

        # Filter posts by the selected tags
        posts = []
        all_posts = storage.all(Post).values()
        for post in all_posts:
            post_tag_ids = [str(tag.id) for tag in post.tags]  # Get the tag IDs for each post
            if any(tag_id in post_tag_ids for tag_id in tag_ids):
                posts.append(post)
    else:
        # No tags selected, show all posts
        posts = storage.all(Post).values()

    tags = storage.all(Tag).values()  # Fetch all tags for display
    return render_template('index.html', posts=posts, tags=tags)
