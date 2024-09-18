#!/usr/bin/python3
"""Defines the Follower class for user following system"""
from sqlalchemy import Column, String, ForeignKey
from models.base_model import BaseModel, Base
from sqlalchemy.orm import relationship

class Follower(BaseModel, Base):
    """Represents a user following another user"""
    __tablename__ = 'followers'

    follower_id = Column(String(60), ForeignKey('users.id'), primary_key=True)
    followed_id = Column(String(60), ForeignKey('users.id'), primary_key=True)

    #Relationship
    follower = relationship('User', foreign_keys=[follower_id], backref='following')
    followed = relationship('User', foreign_keys=[followed_id], backref='followers')
