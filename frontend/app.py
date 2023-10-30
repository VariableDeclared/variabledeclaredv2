from flask import Flask, render_template
import requests
from requests.exceptions import ConnectionError
import os
from argparser import VDotComArgParser
from config.config import (
    VDotComConfigLoader,
    VDotComConfig
)
from blog_app.blog import BlogApp

app = Flask(__name__)

VB_API_URL = os.getenv("VB_API_URL", "https://api.variabledeclared.com")
VDOTCOM_CONFIG: VDotComConfig = VDotComConfigLoader().config

@app.route('/')
def index():
    blog_posts = []
    try:
        response = requests.get(f"{VB_API_URL}/posts")
        blog_posts = response.json()
    except ConnectionError as ex:
        app.logger.warning('Failed to pull posts from the API. Is it up? ex: %s', ex)
    blog_app = VDOTCOM_CONFIG.app
    return render_template('index.html', blog_posts=blog_posts, blog_app=blog_app)

if __name__ == '__main__':
    parser = VDotComArgParser(extra_config=VDOTCOM_CONFIG).parser
    args = parser.parse_args()
    app.run(debug=args.debug, port=args.port)
