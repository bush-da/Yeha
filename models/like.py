#!/usr/bin/python3
"""Defines the Like class"""
from sqlalchemy import Column, String, ForeignKey
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class Like(BaseModel, Base):
    """Represents a 'like' on a post"""
    __tablename__ = 'likes'

    author_id = Column(String(60), ForeignKey('users.id'), primary_key=True)
    post_id = Column(String(60), ForeignKey('posts.id'), primary_key=True)

    # Relationships
    author = relationship("User", back_populates="likes")
    post = relationship("Post", back_populates="likes")
