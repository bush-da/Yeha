#!/usr/bin/python3
"""Defines the Report class for flagging posts"""
from sqlalchemy import Column, String, Enum, Boolean, ForeignKey, CheckConstraint
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class Report(BaseModel, Base):
    """Represents a report for flagging inappropriate posts or comments"""
    __tablename__ = 'reports'

    user_id = Column(String(60), ForeignKey('users.id'), nullable=False)

    post_id = Column(String(60), ForeignKey('posts.id'), nullable=True)
    comment_id = Column(String(60), ForeignKey('comments.id'), nullable=True)

    # Predefined reasons
    reason = Column(Enum('Spam', 'Harassment', 'Misinformation', 'Violence', 'Hate Speech', 'Other', name="flag_reasons"), nullable=False)

    reviewed = Column(Boolean, default=False)  # Track if reviewed
    action_taken = Column(String(64), nullable=True)  # Actions: 'deleted', 'warned', etc.

    # Relationships for navigation, no user relationship for anonymity

    # Report model
    post = relationship('Post', back_populates='reports')
    comment = relationship('Comment', back_populates='reports')

    # Constraint to enforce single target
    __table_args__ = (
        CheckConstraint('(post_id IS NULL) != (comment_id IS NULL)', name='check_single_target'),
    )
