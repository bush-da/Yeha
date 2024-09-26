# from flask import Flask
# from models.engine.db_storage import DBStorage
# from flask_sqlalchemy import SQLAlchemy
# from os import getenv

# app = Flask(__name__)

# # Configure app
# app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql+mysqldb://{}:{}@{}/{}'.format(
#     getenv('YEHA_MYSQL_USER'),
#     getenv('YEHA_MYSQL_PWD'),
#     getenv('YEHA_MYSQL_HOST'),
#     getenv('YEHA_MYSQL_DB')
# )
# app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

# # Initialize storage and db
# storage = DBStorage()
# storage.reload()
# db = SQLAlchemy(app)

# from web_flask import routes
