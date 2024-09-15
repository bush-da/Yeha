#!/usr/bin/python3
"""Association table for posts and tags"""
from sqlalchemy import Column, ForeignKey
from models.base_model import BaseModel, Base

class PostTag(BaseModel, Base):
    __tablename__ = 'post_tags'

    post_id = Column(String(60), ForeignKey('posts.id'), primary_key=True)
    tag_id = Column(String(60), ForeignKey('tags.id'), primary_key=True)
