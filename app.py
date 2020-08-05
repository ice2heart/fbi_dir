import qrcode

import os
import re
from io import BytesIO
from flask import Flask, render_template, send_file, send_from_directory
from environs import Env

env = Env()
env.read_env()
file_path = env.path('CIA_PATH')
qr_ip = env.path('QR_IP')
qr_port = env.int('QR_PORT', 5000)

app = Flask(__name__)


@app.route('/')
def index():
    files = [f for f in os.listdir(file_path) if re.match(r'.*\.cia', f)]
    return render_template('index.html', files=files)


@app.route('/files/<file_name>')
def send_static_file(file_name):
    return send_from_directory(file_path, file_name)


@app.route('/img/<file_name>')
def qr_code(file_name):
    qr = qrcode.QRCode(
        version=1,
        error_correction=qrcode.constants.ERROR_CORRECT_L,
        box_size=10,
        border=4,
    )
    qr.add_data(f'http://{qr_ip}:{qr_port}/files/{file_name}')
    qr.make(fit=True)

    img = qr.make_image(fill_color="black", back_color="white")
    img_io = BytesIO()
    img.save(img_io, 'JPEG', quality=70)
    img_io.seek(0)
    return send_file(img_io, mimetype='image/jpeg')
