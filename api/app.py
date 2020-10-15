import json
import random

from flask import Flask, jsonify, request
from flask_cors import CORS
from error_model import InvalidUsage

application = Flask(__name__)
CORS(application)

JSON_NOMINA = ""
json_verba = ""
json_misc = ""


@application.route('/ping')
def ping_pong():
    return 'pong'


# creates a new quiz word from the nominas
@application.route('/api/v1/nomina/<chapter>', methods=['GET'])
@application.route('/api/v1/nomina', methods=['GET'], defaults={'chapter': None})
def nomina(chapter):
    global JSON_NOMINA
    if JSON_NOMINA == "":
        set_global_lists()

    nomina_json_list = JSON_NOMINA['nomina']

    if chapter is not None:
        chapter_list = list(filter(lambda word_list: word_list['chapter'] == int(chapter), nomina_json_list))
        nomina_json_list = chapter_list

    quiz = create_new_quiz_list(nomina_json_list)

    return jsonify(quiz), 200


# creates a new quiz word from the verbas
@application.route('/api/v1/verba/<chapter>', methods=['GET'])
@application.route('/api/v1/verba', methods=['GET'], defaults={'chapter': None})
def verba(chapter):
    global json_verba
    if json_verba == "":
        set_global_lists()

    verba_json_list = json_verba['verba']

    if chapter is not None:
        chapter_list = list(filter(lambda word_list: word_list['chapter'] == int(chapter), verba_json_list))
        verba_json_list = chapter_list

    quiz = create_new_quiz_list(verba_json_list)

    return jsonify(quiz), 200


# creates a new quiz word from the misc category
@application.route('/api/v1/misc/<chapter>', methods=['GET'])
@application.route('/api/v1/misc', methods=['GET'], defaults={'chapter': None})
def misc(chapter):
    global json_misc
    if json_misc == "":
        set_global_lists()

    misc_json_list = json_misc['misc']

    if chapter is not None:
        chapter_list = \
            list(filter(lambda word_list: word_list['chapter'] == int(chapter), misc_json_list))
        misc_json_list = chapter_list

    quiz = create_new_quiz_list(misc_json_list)

    return jsonify(quiz), 200


@application.route('/api/v1/answer', methods=['POST'])
def check_answer():
    json_body = request.json
    answer = json_body['answer']
    quiz_word = json_body['quizWord']
    category = json_body['categorie']

    correct_answer = False

    local_json_list = ""
    if category == "nomina":
        global json_nomina
        local_json_list = json_nomina['nomina']
    elif category == "verba":
        global json_verba
        local_json_list = json_verba['verba']
    elif category == "misc":
        global json_misc
        local_json_list = json_misc['misc']
    else:
        return jsonify({"error": "please provide a categorie"}), 400

    quiz_answer = list(filter(lambda greek: greek['greek'] == quiz_word, local_json_list))

    if quiz_answer[0]['dutch'] == answer:
        correct_answer = True

    return jsonify({"correctAnswer": correct_answer}), 200


# creates a new quiz word from the nominas
@application.route('/api/v1/chapters', methods=['GET'])
def chapters():
    global JSON_NOMINA
    if JSON_NOMINA == "":
        set_global_lists()

    last_item = JSON_NOMINA['nomina'][-1]
    chapters = last_item['chapter']

    return jsonify({"chapters": chapters}), 200


def set_global_lists():
    global JSON_NOMINA
    global json_verba
    global json_misc
    with open('./api/nomina/wordlist.json') as f:
        JSON_NOMINA = json.load(f)
    with open('./api/verba/wordlist.json') as f:
        json_verba = json.load(f)
    with open('./api/misc/wordlist.json') as f:
        json_misc = json.load(f)
    return


def create_new_quiz_list(json_list):
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
    response = jsonify(error.to_dict())
    response.status_code = error.status_code
    return response


if __name__ == '__main__':
    application.run()
