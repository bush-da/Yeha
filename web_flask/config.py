from os import getenv

class Config:
    SECRET_KEY = getenv('SECRET_KEY', 'your_secret_key')
    SQLALCHEMY_DATABASE_URI = getenv('SQLALCHEMY_DATABASE_URI', 'mysql+mysqldb://{}:{}@{}/{}'.format(
        getenv('YEHA_MYSQL_USER'),
        getenv('YEHA_MYSQL_PWD'),
        getenv('YEHA_MYSQL_HOST'),
        getenv('YEHA_MYSQL_DB')
    ))
    SQLALCHEMY_TRACK_MODIFICATIONS = False
