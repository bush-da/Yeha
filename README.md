
# Yeha Blog Platform

Welcome to **Yeha**, a modern, open-source blogging platform designed to let your voice be heard. Whether you're a passionate writer, a casual blogger, or someone looking to share ideas, Yeha makes it easy to connect with your audience. 🚀

## 🌐 Live Demo and Author Info
- **Live Demo**: [coming soon](https://github.com/bush-da/Yeha)
- **Project Blog**: [Final Project Article](https://www.linkedin.com/posts/daniel-tujuma-702454271_introducing-the-yeha-project-a-web-application-activity-7242165583499132929-qr_j?utm_source=share&utm_medium=member_desktop)
- **Author(s)**: [LinkedIn Profile](https://www.linkedin.com/in/daniel-tujuma-702454271)

---
## 📸Screenshots
Here's a glimpse of the project in action:
## Screenshots

<img src="./assets/screenshot.png" alt="Project Screenshot" width="550" height="300">

## 📖 Introduction

Welcome to Yeha, a clean and simple blog platform built with the power of Python, Flask, and SQLAlchemy! Inspired by a need for expressive and easy-to-use blog websites, Yeha was designed to empower users to publish and share their thoughts, ideas, and creativity.

Its main features include user authentication, blog post creation, tagging, content management, and post engagement through likes. The system allows users to create engaging content, categorize it with tags, and interact with posts through a flexible and intuitive user interface. Whether you're a tech enthusiast or a lifestyle blogger, Yeha is built to empower your creativity.

## 🌟 Inspiration

In the age of information overload, I wanted to create a blog platform that’s simple but powerful, focusing on ease of use for both readers and writers. Whether it’s a personal diary or a professional content platform, Yeha aims to provide a minimalist yet robust experience. The journey was filled with technical challenges—especially in the way we handle relationships between users and posts, as well as ensuring the site can scale as more users join and interact.

## 🎯 Project Goals and Features 

    📝 Dynamic Blog Posting: Create, edit, and delete blog posts with rich text formatting.
    🧑‍💻 User Accounts: Users can sign up, manage their profile, and post content.
    🔖 Tag System: Organize posts by tags for better discoverability.
    ❤️ Like System: Track and display the number of likes for each post.
    🔍 Search: Find posts by tags.
    📊 Admin Dashboard: Monitor posts, tags, and users through a comprehensive admin interface.

## 🛠️ Installation

Getting Yeha up and running on your local environment is simple! Follow these steps:

### Prerequisites
- Python 3.8+
- MySQL
- Flask

### Steps
1. **Clone the Repository**
   ```bash
   git clone https://github.com/bush-da/Yeha.git
   cd Yeha
   ```

2. **Create a Virtual Environment**
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   ```

3. **Install Dependencies**
   ```bash
   pip install -r requirements.txt
   ```

4. **Set up Database**
   Update your `config.py` with your database credentials and run:
   ```bash
   flask db upgrade
   ```

5. **Run the App**
   ```bash
   flask run
   ```

   Now visit `http://127.0.0.1:5000` in your browser.

---

## 🚀 Usage

Yeha is simple yet powerful. Here are the key features you can explore:

- **Create a Blog Post**: Write articles with a full-featured text editor and publish them to the world.
- **Tagging System**: Add tags to your posts, helping readers discover content related to their interests.
- **User Profiles**: Customize your profile with a bio, profile picture, and social links.
- **Post Engagement**: Track likes on posts and follow user activity.
- **Search and Browse**: Browse content by tags, search for articles, and explore trending topics.


## 💡Technical Details 
Database Design:

The database uses a relational model with SQLAlchemy to handle users, posts, tags, and likes. Here's an overview of the key relationships:

    User ↔ Posts: A one-to-many relationship allows users to author multiple posts. Each post tracks its creator using author_id.
    Posts ↔ Tags: A many-to-many relationship where posts can have multiple tags, enabling effective content categorization.
    Posts ↔ Likes: We track the number of likes for each post using a separate table to maintain scalability.
---

## 🤝 Contributing

We love collaboration! If you'd like to contribute to Yeha, please follow these steps:

1. Fork the repository.
2. Create a feature branch (`git checkout -b new-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin new-feature`).
5. Open a Pull Request.

We also encourage issue reporting, feature requests, and suggestions for improving Yeha. Let's build something awesome together!

---

## 🌟 Related Projects

Looking for more open-source projects? Check out these cool projects:

- [Jekyll](https://github.com/jekyll/jekyll) - Simple, blog-aware static site generator for personal, project, or organization sites.
- [Ghost](https://github.com/TryGhost/Ghost) - Professional publishing platform.
- [Flask-Blog](https://github.com/miguelgrinberg/microblog) - A Flask-based blog tutorial project.

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

👋 Thank you for checking out Yeha. We hope you enjoy using the platform as much as we enjoyed building it! Feel free to star the project and share your thoughts. 😊
