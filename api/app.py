"""
Main app module used by flask to serve the application.
If this file continues to grow consider moving part to submodules
"""

import json
import random

from flask import Flask, jsonify, request
from flask_cors import CORS
from error_model import InvalidUsage

application = Flask(__name__)
CORS(application)

JSON_NOMINA = ""
JSON_VERBA = ""
JSON_MISC = ""


@application.route('/ping')
def ping_pong():
    """Endpoint to check if the app is running"""
    return 'pong'


@application.route('/api/v1/nomina/<chapter>', methods=['GET'])
@application.route('/api/v1/nomina', methods=['GET'], defaults={'chapter': None})
def nomina(chapter):
    """Creates a new quiz word from the nomina category"""
    global JSON_NOMINA
    if JSON_NOMINA == "":
        set_global_lists()

    nomina_json_list = JSON_NOMINA['nomina']

    if chapter is not None:
        chapter_list = \
            list(filter(lambda word_list: word_list['chapter'] == int(chapter), nomina_json_list))
        nomina_json_list = chapter_list

    quiz = create_new_quiz_list(nomina_json_list)

    return jsonify(quiz), 200


@application.route('/api/v1/verba/<chapter>', methods=['GET'])
@application.route('/api/v1/verba', methods=['GET'], defaults={'chapter': None})
def verba(chapter):
    """Creates a new quiz word from the verba category"""
    global JSON_VERBA
    if JSON_VERBA == "":
        set_global_lists()

    verba_json_list = JSON_VERBA['verba']

    if chapter is not None:
        chapter_list = \
            list(filter(lambda word_list: word_list['chapter'] == int(chapter), verba_json_list))
        verba_json_list = chapter_list

    quiz = create_new_quiz_list(verba_json_list)

    return jsonify(quiz), 200


@application.route('/api/v1/misc/<chapter>', methods=['GET'])
@application.route('/api/v1/misc', methods=['GET'], defaults={'chapter': None})
def misc(chapter):
    """Creates a new quiz word from the misc category"""
    global JSON_MISC
    if JSON_MISC == "":
        set_global_lists()

    misc_json_list = JSON_MISC['misc']

    if chapter is not None:
        chapter_list = \
            list(filter(lambda word_list: word_list['chapter'] == int(chapter), misc_json_list))
        misc_json_list = chapter_list

    quiz = create_new_quiz_list(misc_json_list)

    return jsonify(quiz), 200


@application.route('/api/v1/answer', methods=['POST'])
def check_answer():
    """Checks the send in answer for correctness"""
    json_body = request.json
    answer = json_body['answer']
    quiz_word = json_body['quizWord']
    category = json_body['categorie']

    correct_answer = False

    local_json_list = ""
    if category == "nomina":
        global JSON_NOMINA
        local_json_list = JSON_NOMINA['nomina']
    elif category == "verba":
        global JSON_VERBA
        local_json_list = JSON_VERBA['verba']
    elif category == "misc":
        global JSON_MISC
        local_json_list = JSON_MISC['misc']
    else:
        return jsonify({"error": "please provide a categorie"}), 400

    quiz_answer = list(filter(lambda greek: greek['greek'] == quiz_word, local_json_list))

    if quiz_answer[0]['dutch'] == answer:
        correct_answer = True

    return jsonify({"correctAnswer": correct_answer}), 200


@application.route('/api/v1/chapters', methods=['GET'])
def chapters():
    """Returns the latest chapter number"""
    global JSON_NOMINA
    if JSON_NOMINA == "":
        set_global_lists()

    last_item = JSON_NOMINA['nomina'][-1]
    last_chapter = last_item['chapter']

    return jsonify({"chapters": last_chapter}), 200


def set_global_lists():
    """Sets the global lists to the contents of the different json files"""
    global JSON_NOMINA
    global JSON_VERBA
    global JSON_MISC
    with open('./api/nomina/wordlist.json') as nomina_list:
        JSON_NOMINA = json.load(nomina_list)
    with open('./api/verba/wordlist.json') as verba_list:
        JSON_VERBA = json.load(verba_list)
    with open('./api/misc/wordlist.json') as misc_list:
        JSON_MISC = json.load(misc_list)


def create_new_quiz_list(json_list):
    """Creates a new quiz list for the user
    it will return 5 different answers from the json_list
    if the number of answers is lower than 5 than the length of the json_list + 1 is used
    """
    quiz = []

    random_entry = random.choice(list(json_list))
    answer = random_entry['dutch']
    quiz_word = random_entry['greek']

    quiz.append(quiz_word)
    quiz.append(answer)

    number_of_answers = 5
    if len(json_list) < number_of_answers:
        number_of_answers = len(json_list) + 1

    while len(quiz) != number_of_answers:
        rand_entry = random.choice(list(json_list))
        if rand_entry['dutch'] not in quiz:
            quiz.append(rand_entry['dutch'])
    return quiz


@application.errorhandler(InvalidUsage)
def handle_invalid_usage(error):
    """Error handler for the api"""
    response = jsonify(error.to_dict())
    response.status_code = error.status_code
    return response


if __name__ == '__main__':
    application.run()
