from flask import Flask, render_template
import requests
import os

app = Flask(__name__)

VB_API_URL = os.getenv("VB_API_URL", "https://api.variabledeclared.com")
@app.route('/')
def index():
    # Add logic to retrieve blog posts from a database or file
    # For now, let's assume a list of sample blog posts
    

    blog_posts = []

    response = requests.get(f"{VB_API_URL}/posts")
    blog_posts = response.json()
    return render_template('index.html', blog_posts=blog_posts)

if __name__ == '__main__':
    app.run(debug=True, port=8080)