#!/usr/bin/python3
"""Defines the Content class"""
from sqlalchemy import Column, Text, ForeignKey
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class Content(BaseModel, Base):
    """Content attributes that will be posted on the website"""
    __tablename__ = 'contents'

    file_url = Column(Text, nullable=False)
    post_id = Column(String(60), ForeignKey('posts.id'), nullable=False)

    # Relationship
    post = relationship("Post", back_populates="contents")
