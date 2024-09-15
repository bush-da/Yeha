#!/usr/bin/python3
"""Defines the User class"""
from sqlalchemy import Column, String, Text, Enum, Boolean
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class User(BaseModel, Base):
    """Represents a user in the Yeha application"""
    __tablename__ = 'users'

    username = Column(String(45), nullable=False)
    email = Column(String(45), nullable=False, unique=True)
    password = Column(String(128), nullable=False)
    profile_picture = Column(String(255), nullable=True)
    gender = Column(Enum('male', 'female', 'other'))
    bio = Column(Text, nullable=True)
    location = Column(String(100), nullable=True)
    website = Column(String(100), nullable=True)
    is_admin = Column(Boolean, default=False)

    # Relationships
    posts = relationship("Post", back_populates="author", cascade="all, delete-orphan")
    comments = relationship("Comment", back_populates="author", cascade="all, delete-orphan")
    likes = relationship("Like", back_populates="author", cascade="all, delete-orphan")
