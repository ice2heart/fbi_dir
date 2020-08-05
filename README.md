## Simple app for remote install CIA files

required some env

| env | desctiption | example |
| --- | --- | --- |
| CIA_PATH | path to yours cia files | ~/Downloads  |
| QR_IP | ip for qr codes | 10.74.74.23 |
| QR_PORT | port for qr codes (default 5000) | 5000 |

## installation 

1. install with pipenv

```pipenv install```

2. install with pip

``` pip install -r requirements.txt```

## How to run

1. run with pipenv

```CIA_PATH=~/Downloads QR_IP=10.74.74.23 FLASK_APP=app.py pipenv run flask run --host 0.0.0.0 --with-threads```

2. run directly 

```CIA_PATH=~/Downloads QR_IP=10.74.74.23 FLASK_APP=app.py flask run --host 0.0.0.0 --with-threads```