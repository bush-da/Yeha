from flask import Blueprint, render_template

about_bp = Blueprint('about', __name__, url_prefix='/about')

@about_bp.route('/', methods=['GET'])
def about():
    return render_template('about.html')
