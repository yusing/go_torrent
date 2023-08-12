import flask

app = flask.Flask(__name__)

@app.route('/', methods=['POST'])
def error_receiver():
    msg = flask.request.get_data(as_text=True)
    print(msg)
    return 'OK'

# host on port 5999
app.run(host='localhost', port=5999)
