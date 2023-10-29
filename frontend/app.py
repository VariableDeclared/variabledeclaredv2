from flask import Flask, render_template
import requests
import os
from argparser import VDotComArgParser
from config.config import (
    VDotComConfigLoader,
    VDotComConfig
)
from requests.exceptions import ConnectionError

app = Flask(__name__)

VB_API_URL = os.getenv("VB_API_URL", "https://api.variabledeclared.com")
@app.route('/')
def index():
    blog_posts = []
    try:
        response = requests.get(f"{VB_API_URL}/posts")
        blog_posts = response.json()
    except ConnectionError as ex:
        app.logger.warning('Failed to pull posts from the API. Is it up? ex: %s', ex)

    return render_template('index.html', blog_posts=blog_posts)

if __name__ == '__main__':
    vdotcom_config: VDotComConfig = VDotComConfigLoader().config
    parser = VDotComArgParser(extra_config=vdotcom_config).parser
    args = parser.parse_args()
    app.run(debug=args.debug, port=args.port)
