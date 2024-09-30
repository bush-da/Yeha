from flask import Blueprint, render_template, request, redirect, url_for, session
from models.post import Post
from models.tag import Tag
from models.follower import Follower
from models import storage

home_bp = Blueprint('home', __name__)

@home_bp.route('/', methods=['GET'])
def redirect_home():
    return redirect(url_for('home.index'))

@home_bp.route('/home', methods=['GET'])
def index():
    logged_in = 'user_id' in session
    tag_ids = request.args.get('tags')

    # Fetch all posts and tags from storage
    posts = list(storage.all(Post).values())
    tags = list(storage.all(Tag).values())

    if tag_ids:
        # Split the comma-separated tag IDs and convert to integers for comparison
        tag_ids = set(tag_ids.split(','))  # Using a set for efficient lookup

        # Filter posts that contain any of the selected tags
        filtered_posts = []
        for post in posts:
            post_tag_ids = {str(tag.id) for tag in post.tags}  # Post tag IDs as a set of strings
            if post_tag_ids.intersection(tag_ids):  # Check if there's any match
                filtered_posts.append(post)

        posts = filtered_posts  # Update the posts list with filtered posts

    follower_map = {}
    for post in posts:
        post.contents = sorted(post.contents, key=lambda c: c.paragraph)
        followers = storage.all(Follower).values()
        # Filter followers for the current post's author
        follower_map[post.author.id] = [follower.follower_id for follower in followers if follower.followed_id == post.author.id]

    return render_template('index.html', posts=posts, tags=tags, follower_map=follower_map, logged_in=logged_in)
