#!/usr/bin/python3
"""Association table for posts and tags"""
from sqlalchemy import Column, ForeignKey, String, Table
from models.base_model import BaseModel, Base

post_tag = Table(
    'post_tags',
    Base.metadata,
    Column('post_id', String(60), ForeignKey('posts.id'), primary_key=True),
    Column('tag_id', String(60), ForeignKey('tags.id'), primary_key=True),
)
